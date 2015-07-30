package netlink

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

type Listener struct {
	sock        *Socket
	Messagechan chan Message
	sendqueue   chan *Message
	currSeq     uint32
	lock        sync.Mutex
	socketLock  sync.Mutex
	echan       chan error
}

// Used as an atomic counter for sequence numbering.
// No check is made to see that sequences aren't still in use on roll-over.
func (listener *Listener) Seq() (out uint32) {
	listener.lock.Lock()
	listener.currSeq++
	out = listener.currSeq
	listener.lock.Unlock()
	return
}

func (listener *Listener) Close() {
	//close(listener.Messagechan)
	listener.sock.Close()
}

// Send a message.  If SequenceNumber is unset, Seq() will be used
// to generate one.
func (listener *Listener) Query(msg *Message) (err error) {
	listener.sendqueue <- msg
	return
}

func NewListener(nlfm NetlinkFamily) (listener *Listener, err error) {
	mysock, err := Dial(nlfm)
	listener = nil
	if err != nil {
		log.Panicf("Can't dial netlink socket: %v", err)
		err = err
		return
	}

	listener = &Listener{sock: mysock, Messagechan: make(chan Message, 10), sendqueue: make(chan *Message, 10), currSeq: 0}
	return
}

func (listener *Listener) startListening() {
	r := bufio.NewReader(listener.sock)
	for listener.sock.IsOpen() {
		peekedBytes, err := r.Peek(1)
		if peekedBytes == nil {
			log.Println("Didn't peeked any bytes")
			continue
		}
		if err != nil && err == bufio.ErrNegativeCount {
			// Most probably the socket is closed
			log.Printf("Netlink socket seems to be closed, checking if socket is open")
			continue
		}
		//listener.socketLock.Lock()
		messages := make([]*Message, 0, 10)
		msg, err := ReadMessage(r)
		messages = append(messages, msg)
		for !isLastMessage(msg) {
			msg, err = ReadMessage(r)
			messages = append(messages, msg)
		}
		if msg.Header.MessageSequence() == listener.currSeq {
			// We received the response to a request, so we can send the next Query
			listener.socketLock.Unlock()
		}
		//listener.socketLock.Unlock()
		for i := range messages {
			msg := messages[i]
			listener.handleMessage(msg)
		}
		runtime.Gosched()
	}
}

func (listener *Listener) sendError(err error) {
	if listener.echan != nil {
		listener.echan <- err
	} else {
		log.Fatalf("Can't parse netlink message: %v", err)
	}
}

func (listener *Listener) handleMessage(msg *Message) {
	if msg != nil {
		if msg.Header.MessageType() == NLMSG_ERROR {
			errmsg := &Error{}
			err := errmsg.UnmarshalNetlink(msg.Body)
			if err != nil {
				log.Panicf("Can't unmarshall netlink error message: %v", err)
			} else {
				err = errors.New(fmt.Sprintf("Netlink Error (%d) for message with sequence (%d): %s. Header %v Body: %v", errmsg.Code(), msg.Header.MessageSequence(), errmsg.Error(), msg.Header, msg.Body))
				listener.echan <- err
			}

		} else {
			listener.Messagechan <- *msg
		}

	} else {
		log.Fatalf("Netlink message was null")
	}
}

func (listener *Listener) startWriting() {
	for listener.sock.IsOpen() {
		select {
		case msg := <-listener.sendqueue:
			listener.socketLock.Lock()
			if msg.Header.MessageSequence() == 0 {
				msg.Header.SetMessageSequence(listener.Seq())
			}
			ob, err := msg.MarshalNetlink()
			if err == nil {
				_, err = listener.sock.Write(ob)
				runtime.Gosched()
			}
		default:
			runtime.Gosched()
		}
	}
}

func (listener *Listener) Start(echan chan error) (err error) {
	// ^uint32 is MAX UNIT and means that we want to listen to all multicast groups
	listener.echan = echan
	err = listener.sock.Bind(uint32(os.Getpid()), ^uint32(0))
	if err != nil {
		log.Panicf("Can't bind to netlink socket: %v", err)
		err = err
		return
	}
	go listener.startWriting()
	go listener.startListening()
	return
}

func isLastMessage(msg *Message) bool {
	if msg == nil {
		return true
	}
	if msg.Header.MessageType() == NLMSG_DONE {
		return true
	}
	return msg.Header.MessageFlags()&NLM_F_MULTI == 0
}

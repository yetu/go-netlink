package netlink

import (
	"bufio"
	"log"
	"os"
	"sync"
)

type Listener struct {
	sock        *Socket
	Messagechan chan Message
	nextSeq     uint32
	lock        sync.Mutex
}

// Used as an atomic counter for sequence numbering.
// No check is made to see that sequences aren't still in use on roll-over.
func (listener *Listener) Seq() (out uint32) {
	listener.lock.Lock()
	out = listener.nextSeq
	listener.nextSeq++
	listener.lock.Unlock()
	return
}

func (listener *Listener) Close() {
	//close(listener.Messagechan)
	listener.sock.Close()
}

// Send a message.  If SequenceNumber is unset, Seq() will be used
// to generate one.
func (listener *Listener) Query(msg Message, l int) (ch chan Message, err error) {
	if msg.Header.MessageSequence() == 0 {
		msg.Header.SetMessageSequence(listener.Seq())
	}
	ob, err := msg.MarshalNetlink()
	if err == nil {
		_, err = listener.sock.Write(ob)
	}
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

	listener = &Listener{sock: mysock, Messagechan: make(chan Message), nextSeq: 1}
	return
}

func (listener *Listener) Start(echan chan error) (err error) {
	// ^uint32 is MAX UNIT and means that we want to listen to all multicast groups
	err = listener.sock.Bind(uint32(os.Getpid()), ^uint32(0))
	if err != nil {
		log.Panicf("Can't bind to netlink socket: %v", err)
		err = err
		return
	}
	r := bufio.NewReader(listener.sock)
	for listener.sock.IsOpen() {
		_, err = r.Peek(1)
		if err != nil && err == bufio.ErrNegativeCount {
			// Most probably the socket is closed
			continue
		}
		msg, err := ReadMessage(r)
		if err != nil {
			if echan != nil {
				echan <- err
			} else {
				log.Fatalf("Can't parse netlink message: %v", err)
			}
		} else if msg != nil {
			listener.Messagechan <- *msg
		} else {
			log.Fatalf("Netlink message was null")
		}
	}
	return
}

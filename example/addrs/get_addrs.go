package main

/* The output of this utility is JSON, but not intended to be easily human-readable.
   At this moment it exists to test the rtnetlink/route subsystem */

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "github.com/yetu/go-netlink/rtnetlink/addr"
import "github.com/yetu/go-netlink/rtnetlink"
import "log"
import "github.com/yetu/go-netlink"

func logec(c chan error) {
	for i := range c {
		log.Printf("Error: %v", i)
	}
}

func main() {
	nlmsg, err := netlink.NewMessage(rtnetlink.RTM_GETADDR, netlink.NLM_F_DUMP|netlink.NLM_F_REQUEST, &addr.Header{})
	if err != nil {
		log.Panicf("Couldn't construct message: %v", err)
	}
	//log.Printf("Dialing: %v", nlmsg)
	nlsock, err := netlink.Dial(netlink.NETLINK_ROUTE)
	if err != nil {
		log.Panicf("Couldn't dial netlink: %v", err)
	}
	h := netlink.NewHandler(nlsock)
	ec := make(chan error)
	go logec(ec)
	go h.Start(ec)
	//log.Printf("Sending query: %v", nlmsg)
	c, err := h.Query(*nlmsg, 1)
	//log.Printf("Sent query: %v", nlmsg.Header)
	if err != nil {
		log.Panicf("Couldn't write netlink: %v", err)
	}
	for i := range c {
		if i.Header.MessageType() == netlink.NLMSG_DONE {
			break
		}
		switch i.Header.MessageType() {
		case rtnetlink.RTM_NEWADDR:
			hdr := &addr.Header{}
			msg := rtnetlink.NewMessage(hdr, nil)
			err = msg.UnmarshalNetlink(i.Body)
			if err == nil {
				log.Printf("Family: %s; Length: %d; Flags: %v; Scope: %v; IFIndex: %d",
					hdr.AddressFamily(), hdr.PrefixLength(), hdr.Flags(), hdr.Scope(),
					hdr.InterfaceIndex())

				for i := range msg.Attributes {
					attr := msg.Attributes[i]
					log.Printf("Attribute[%s] = %v", addr.AttributeTypeStrings[attr.Type], attr.Body)
				}
			} else {
				log.Printf("Unmarshal error: %v", err)
			}
		default:
			log.Printf("Unknown type: %v", i)
		}
		if err != nil {
			log.Printf("Handler error: %v", err)
		}
	}
}

package main

/* The output of this utility is JSON, but not intended to be easily human-readable.
   At this moment it exists to test the rtnetlink/link subsystem */

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "github.com/yetu/go-netlink/rtnetlink/link"
import "github.com/yetu/go-netlink/rtnetlink"
import "github.com/yetu/go-netlink/xfrm"
import "log"
import "github.com/yetu/go-netlink"

func main() {
	nlmsg, err := netlink.NewMessage(xfrm.XFRM_MSG_GETPOLICY, netlink.NLM_F_DUMP|netlink.NLM_F_REQUEST, &link.Header{})
	if err != nil {
		log.Panicf("Couldn't construct message: %v", err)
	}
	nlsock, err := netlink.Dial(netlink.NETLINK_XFRM)
	if err != nil {
		log.Panicf("Couldn't dial netlink: %v", err)
	}
	h := netlink.NewHandler(nlsock)
	ec := make(chan error)
	go h.Start(ec)
	c, err := h.Query(*nlmsg, 1)
	if err != nil {
		log.Panicf("Couldn't write netlink: %v", err)
	}
	for i := range c {
		if i.Header.MessageType() == netlink.NLMSG_DONE {
			break
		}
		if i.Header.MessageType() == netlink.NLMSG_ERROR {
			emsg := &netlink.Error{}
			err = emsg.UnmarshalNetlink(i.Body)
			if err == nil && emsg.Code() != 0 {
				log.Printf("Netlink error: %v", emsg.Error())
			}
			break
		}
		switch i.Header.MessageType() {
		case xfrm.XFRM_MSG_NEWPOLICY:
			hdr := &xfrm.Header{}
			msg := rtnetlink.NewMessage(hdr, nil)
			err = msg.UnmarshalNetlink(i.Body)
			if err == nil {
				log.Printf("Policy: %v", hdr)
				for i := range msg.Attributes {
					log.Printf("Attribute[%d]: %v", i, msg.Attributes[i])
				}
			} else {
				log.Printf("Unmarshal error: %v", err)
			}
		default:
			log.Printf("Unknown type: %v", i)
		}
	}
}

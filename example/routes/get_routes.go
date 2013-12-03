package main

/* The output of this utility is JSON, but not intended to be easily human-readable.
   At this moment it exists to test the rtnetlink/route subsystem */

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "bitbucket.org/vtolstov/go-netlink/rtnetlink/link"
import "bitbucket.org/vtolstov/go-netlink/rtnetlink/route"
import "bitbucket.org/vtolstov/go-netlink/rtnetlink"
import "log"
import "bitbucket.org/vtolstov/go-netlink"
import "encoding/binary"

func main() {
	nlmsg, err := netlink.NewMessage(rtnetlink.RTM_GETROUTE, netlink.NLM_F_DUMP|netlink.NLM_F_REQUEST, &link.Header)
	if err != nil {
		log.Panicf("Couldn't construct message: %v", err)
	}
	nlsock, err := netlink.Dial(netlink.NETLINK_ROUTE)
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
		switch i.Header.MessageType() {
		case rtnetlink.RTM_NEWROUTE:
			hdr := &route.Header{}
			msg := rtnetlink.NewMessage(hdr, nil)
			err = msg.UnmarshalNetlink(i.Body)

			if err == nil {
				log.Printf("Route: %v (%d/%d) TOS: %d; (Table: %v; Origin: %v; Scope: %v; Type: %v; Flags: %v",
					hdr.AddressFamily(), hdr.AddressDestLength(), hdr.AddressSourceLength(),
					hdr.TOS(), hdr.RoutingTable(), hdr.RouteOrigin(), hdr.AddressScope(),
					hdr.RouteType(), hdr.Flags())
				for i := range msg.Attributes {
					log.Printf("Attribute[%d]: (%v) %v", i, route.AttributeTypeStrings[msg.Attributes[i].Type], msg.Attributes[i])
				}
			} else {
				log.Printf("Unmarshal error: %v", err)
			}
		default:
			log.Printf("Unknown type: %v", i)
		}
	}
}

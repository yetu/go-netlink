package xfrm

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "github.com/yetu/go-netlink"
import "github.com/yetu/go-netlink/rtnetlink"

import (
	//"encoding/binary"
	"errors"
)

const HEADER_LENGTH = 168

type Header [HEADER_LENGTH]byte

func NewHeader(fam rtnetlink.Family) (hdr *Header) {
	hdr = &Header{byte(fam)}
	return
}

func (self Header) Len() int { return HEADER_LENGTH }
func (self *Header) UnmarshalNetlink(in []byte) (err error) {
	if len(in) != HEADER_LENGTH {
		err = errors.New("Wrong length for Header")
	} else {
		copy(self[0:HEADER_LENGTH], in[0:HEADER_LENGTH])
	}
	return
}

func (self Header) MarshalNetlink() (out []byte, err error) {
	out = netlink.Padded(self[0:HEADER_LENGTH])
	return
}

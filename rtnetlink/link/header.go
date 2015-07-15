package link

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "github.com/yetu/go-netlink"
import "github.com/yetu/go-netlink/rtnetlink"

import (
	"encoding/binary"
	"errors"
)

const HEADER_LENGTH = 16

type Header [16]byte

func NewHeader(fam rtnetlink.Family, itype uint16, iindex uint32, flags, changes Flags) (hdr *Header) {
	hdr = &Header{byte(fam)}
	binary.LittleEndian.PutUint16(hdr[2:4], itype)
	binary.LittleEndian.PutUint32(hdr[4:8], iindex)
	binary.LittleEndian.PutUint32(hdr[8:12], uint32(flags))
	binary.LittleEndian.PutUint32(hdr[12:16], uint32(changes))
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

func ParseMessage(msg netlink.Message) (ret *rtnetlink.Message, err error) {
	hdr := &Header{}
	ret = rtnetlink.NewMessage(hdr, nil)
	err = ret.UnmarshalNetlink(msg.Body)
	if err != nil {
		ret = nil
		return
	}
	return
}

func (self Header) Flags() Flags                      { return Flags(binary.LittleEndian.Uint32(self[8:12])) }
func (self *Header) SetFlags(f Flags)                 { binary.LittleEndian.PutUint32(self[8:12], uint32(f)) }
func (self Header) InterfaceFamily() rtnetlink.Family { return rtnetlink.Family(self[0]) }
func (self Header) InterfaceType() uint16             { return binary.LittleEndian.Uint16(self[2:4]) }
func (self Header) InterfaceIndex() uint32            { return binary.LittleEndian.Uint32(self[4:8]) }
func (self Header) InterfaceChanges() Flags           { return Flags(binary.LittleEndian.Uint32(self[12:16])) }

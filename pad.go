package netlink

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Copyright (c) 2013, Vishvananda Ishaya. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/
const (
	ALIGN_MINUS_ONE     = 3
	NOT_ALIGN_MINUS_ONE = ^3
)

// Returns a padded version of bytes
func Padded(in []byte) (out []byte) {
	size := len(in)
	align := Align(size)
	if align != size {
		out = make([]byte, align)
		copy(out, in)
	} else {
		out = in
	}
	return
}

// Returns a four byte allined value.
func Align(pos int) int {
	return (pos + ALIGN_MINUS_ONE) & NOT_ALIGN_MINUS_ONE
}

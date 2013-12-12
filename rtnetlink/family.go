package rtnetlink

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "syscall"

// Family is already defined in syscall, but we re-define
// them here for type-safety.
type Family uint8

const (
	AF_UNSPEC  Family = syscall.AF_UNSPEC
	AF_INET    Family = syscall.AF_INET
	AF_INET6   Family = syscall.AF_INET6
	AF_BRIDGE  Family = syscall.AF_BRIDGE
)

// Returns the String representation of the Family,
// or "" if the family name is unknown.
func (self Family) String() (out string) {
	switch self {
	default:
		out = ""
	case AF_UNSPEC:
		out = "AF_UNSPEC"
	case AF_INET:
		out = "AF_INET"
	case AF_INET6:
		out = "AF_INET6"
	case AF_BRIDGE:
		out = "AF_BRIDGE"
	}
	return
}

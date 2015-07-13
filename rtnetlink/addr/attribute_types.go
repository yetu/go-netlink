package addr

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import "github.com/yetu/go-netlink"

const (
	IFA_UNSPEC netlink.AttributeType = iota
	IFA_ADDRESS
	IFA_LOCAL
	IFA_LABEL
	IFA_BROADCAST
	IFA_ANYCAST
	IFA_CACHEINFO
	IFA_MULTICAST
	IFA_MAX = IFA_MULTICAST
)

var AttributeTypeStrings = map[netlink.AttributeType]string{
	IFA_UNSPEC:    "IFA_UNSPEC",
	IFA_ADDRESS:   "IFA_ADDRESS",
	IFA_LOCAL:     "IFA_LOCAL",
	IFA_LABEL:     "IFA_LABEL",
	IFA_BROADCAST: "IFA_BROADCAST",
	IFA_ANYCAST:   "IFA_ANYCAST",
	IFA_CACHEINFO: "IFA_CACHEINFO",
	IFA_MULTICAST: "IFA_MULTICAST",
}

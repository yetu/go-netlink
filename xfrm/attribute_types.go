package xfrm

/*
  Copyright (c) 2013, Vishvananda Ishaya. All rights reserved.
  Original Author: Vishvananda Ishaya <vishvananda@gmail.com>

  See LICENSE for details
*/

import "github.com/vishvananda/go-netlink"

const (
	XFRMA_UNSPEC netlink.AttributeType = iota
	XFRMA_ALG_AUTH
	XFRMA_ALG_CRYPT
	XFRMA_ALG_COMP
	XFRMA_ENCAP
	XFRMA_TMPL
	XFRMA_SA
	XFRMA_POLICY
	XFRMA_SEC_CTX
	XFRMA_LTIME_VAL
	XFRMA_REPLAY_VAL
	XFRMA_REPLAY_THRESH
	XFRMA_ETIMER_THRESH
	XFRMA_SRCADDR
	XFRMA_COADDR
	XFRMA_LASTUSED
	XFRMA_POLICY_TYPE
	XFRMA_MIGRATE
	XFRMA_ALG_AEAD
	XFRMA_KMADDRESS

	XFRMA_MAX = XFRMA_KMADDRESS
)

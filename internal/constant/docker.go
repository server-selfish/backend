package constant

import "net/netip"

var (
	ALL_ADDR, _     = netip.ParseAddr("0.0.0.0")
	MAXIMUM_RESTART = 999
)

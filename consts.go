package goosr

// Needs for parsing Windows Ticks
const epoch = 621355968000000000

const (
	OSU = iota
	TAIKO
	CTB
	MANIA
)

const (
	MOUSELEFT = 1 << iota
	MOUSERIGHT
	K1
	K2
	SMOKE
)

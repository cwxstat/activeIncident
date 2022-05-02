package constants

import (
	"time"
)

var (
	WebCadChester     = "https://webcad.chesco.org/WebCad/webcad.asp"
	WebCadMontcoPrint = "https://webapp07.montcopa.org/eoc/cadinfo/livecad.asp?print=yes"
	WebCadMontco      = "https://webapp07.montcopa.org/eoc/cadinfo/"
	RefreshRate       = time.Second * 70
	ErrorBackoff      = time.Second * 200
	MontcoZipCodes    = []int{
		19027,
		18041,
		18426,
		18964,
		19044,
		19454,
		19428}
)

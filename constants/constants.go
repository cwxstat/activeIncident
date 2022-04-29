package constants

import (
	"time"
)

var (
	WebCadChester = "https://webcad.chesco.org/WebCad/webcad.asp"
	WebCadMontco  = "https://webapp07.montcopa.org/eoc/cadinfo/livecad.asp?print=yes"
	RefreshRate   = time.Second * 70
	ErrorBackoff  = time.Second * 200
)

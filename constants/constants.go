package constants

import (
	"time"
)

var (
	WebCadURL    = "https://webapp07.montcopa.org/eoc/cadinfo/"
	RefreshRate  = time.Second * 70
	ErrorBackoff = time.Second * 200
)

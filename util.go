package dcodec

import (
	"github.com/deepch/goav/avformat"
)

func init() {
	avformat.AvRegisterAll()
}

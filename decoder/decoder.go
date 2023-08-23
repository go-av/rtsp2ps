package decoder

import (
	"time"

	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/pion/rtp"
)

type Decoder interface {
	Decode(pkt *rtp.Packet) ([][]byte, time.Duration, error)
	Format() formats.Format
	CheckError(err error) bool
}

package decoder

import (
	"time"

	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/formats/rtph264"
	"github.com/pion/rtp"
)

func NewH264Decoder(f *formats.H264) (Decoder, error) {
	decoder, err := f.CreateDecoder2()
	if err != nil {
		return nil, err
	}

	return &H264Decoder{
		format:  f,
		decoder: decoder,
	}, nil
}

type H264Decoder struct {
	format  *formats.H264
	decoder *rtph264.Decoder
}

func (p *H264Decoder) Decode(pkt *rtp.Packet) ([][]byte, time.Duration, error) {
	return p.decoder.Decode(pkt)
}

func (p *H264Decoder) Format() formats.Format {
	return p.format
}

func (p *H264Decoder) CheckError(err error) bool {
	if err == rtph264.ErrNonStartingPacketAndNoPrevious || err == rtph264.ErrMorePacketsNeeded {
		return false
	}
	return true
}

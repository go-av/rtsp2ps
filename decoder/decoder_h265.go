package decoder

import (
	"time"

	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/formats/rtph265"
	"github.com/pion/rtp"
)

func NewH265Decoder(f *formats.H265) (Decoder, error) {
	decoder, err := f.CreateDecoder2()
	if err != nil {
		return nil, err
	}

	return &H265Decoder{
		format:  f,
		decoder: decoder,
	}, nil
}

type H265Decoder struct {
	format  *formats.H265
	decoder *rtph265.Decoder
}

func (p *H265Decoder) Decode(pkt *rtp.Packet) ([][]byte, time.Duration, error) {
	return p.decoder.Decode(pkt)
}

func (p *H265Decoder) Format() formats.Format {
	return p.format
}

func (p *H265Decoder) CheckError(err error) bool {
	if err == rtph265.ErrNonStartingPacketAndNoPrevious || err == rtph265.ErrMorePacketsNeeded {
		return false
	}
	return true
}

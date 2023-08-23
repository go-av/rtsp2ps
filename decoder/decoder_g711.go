package decoder

import (
	"time"

	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/formats/rtpsimpleaudio"

	"github.com/pion/rtp"
)

func NewG711Decoder(f *formats.G711) (Decoder, error) {
	decoder, err := f.CreateDecoder2()
	if err != nil {
		return nil, err
	}

	return &G711Decoder{
		format:  f,
		decoder: decoder,
	}, nil
}

type G711Decoder struct {
	format  *formats.G711
	decoder *rtpsimpleaudio.Decoder
}

func (p *G711Decoder) Decode(pkt *rtp.Packet) ([][]byte, time.Duration, error) {
	payload, pts, err := p.decoder.Decode(pkt)
	return [][]byte{payload}, pts, err
}

func (p *G711Decoder) Format() formats.Format {
	return p.format
}

func (p *G711Decoder) CheckError(err error) bool {
	return true
}

package ps

import (
	"fmt"
	"time"

	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/pion/rtp"
	"github.com/sirupsen/logrus"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

func NewPSMuxer() *PSMuxer {
	return &PSMuxer{
		mux:     mpeg2.NewPsMuxer(),
		streams: make(map[string]uint8),
	}
}

type PSMuxer struct {
	mux     *mpeg2.PSMuxer
	streams map[string]uint8
}

func (m *PSMuxer) OnPacket(onPacket func([]byte)) {
	m.mux.OnPacket = onPacket
}

func (m *PSMuxer) Muxer(format formats.Format, nalus [][]byte, dts time.Duration) error {
	streamID, ok := m.streams[format.Codec()]
	if !ok || streamID == 0 {
		switch format.(type) {
		case *formats.H264:
			streamID = m.mux.AddStream(mpeg2.PS_STREAM_H264)
		case *formats.H265:
			streamID = m.mux.AddStream(mpeg2.PS_STREAM_H265)
		case *formats.G711:
			if format.PayloadType() == 0 {
				streamID = m.mux.AddStream(mpeg2.PS_STREAM_G711U)
			} else {
				streamID = m.mux.AddStream(mpeg2.PS_STREAM_G711A)
			}
		}
		m.streams[format.Codec()] = streamID
		fmt.Println("streamID", streamID, format.Codec(), format.ClockRate(), dts.Milliseconds())
	}

	if streamID == 0 {
		logrus.Debug("streamID is 0 ", format.Codec(), format.ClockRate())
		return nil
	}

	for _, nalu := range nalus {
		err := m.mux.Write(streamID, append([]byte{00, 00, 00, 01}, nalu...), uint64(dts.Milliseconds()), uint64(dts.Milliseconds()))
		if err != nil {
			logrus.Errorf("ps muxer error: %s", err)
		}
	}

	return nil
}

func (m *PSMuxer) ClockRate() uint32 {
	return 90000
}

func (m *PSMuxer) Payloader() rtp.Payloader {
	return &PSPayloader{}
}

func (m *PSMuxer) PlayloadType() uint8 {
	return 96
}

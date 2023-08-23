package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/media"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/go-av/rtsp2ps/decoder"
	"github.com/go-av/rtsp2ps/ps"
	"github.com/pion/rtp"
	"github.com/sirupsen/logrus"
)

func main() {
	rtspURL := flag.String("rtsp-url", "rtsp://", "RTSP URL")
	outputDir := flag.String("output-dir", "./output", "输出目录")
	flag.Parse()
	client := &gortsplib.Client{}
	client.Transport = ptr(gortsplib.TransportTCP)
	u, err := url.Parse(*rtspURL)
	if err != nil {
		panic(err)
	}

	if err := client.Start(u.Scheme, u.Host); err != nil {
		panic(err)
	}

	medias, baseURL, _, err := client.Describe(u)
	if err != nil {
		panic(err)
	}

	if len(medias) < 1 {
		panic(fmt.Errorf("invalid medias %v", medias))
	}

	for _, m := range medias {
		if _, err := client.Setup(m, baseURL, 0, 0); err != nil {
			panic("play rtsp failed")
		}
	}

	name := time.Now().Format("20060102150405")
	os.MkdirAll(*outputDir, os.ModePerm)
	file, _ := os.Create(*outputDir + "/" + name + ".ps")
	h264file, _ := os.Create(*outputDir + "/" + name + ".h264")

	decoders := map[string]decoder.Decoder{}
	muxer := ps.NewPSMuxer()

	muxer.OnPacket(func(packet []byte) {
		file.Write(packet)
	})

	// RTP
	// packetizer := rtp.NewPacketizer(1450, muxer.PlayloadType(), 1111, muxer.Payloader(), rtp.NewRandomSequencer(), muxer.ClockRate())
	// muxer.OnPacket(func(packet []byte) {
	// 	for _, pkt := range packetizer.Packetize(packet, 90000) {
	// 		file.Write(pkt.Payload)
	// 	}
	// })

	client.OnPacketRTPAny(func(m *media.Media, f formats.Format, packet *rtp.Packet) {
		decode, ok := decoders[f.Codec()]
		if !ok {
			switch x := f.(type) {
			case *formats.H264:
				decode, err = decoder.NewH264Decoder(x)
				if err == nil {
					decoders[f.Codec()] = decode
				}
			case *formats.H265:
				decode, err = decoder.NewH265Decoder(x)
				if err == nil {
					decoders[f.Codec()] = decode
				}
			case *formats.G711:
				decode, err = decoder.NewG711Decoder(x)
				if err == nil {
					decoders[f.Codec()] = decode
				}
			default:
				return
			}
			if err != nil {
				logrus.Error(err)
				return
			}
		}

		if decode == nil {
			return
		}

		nalus, dts, err := decode.Decode(packet)
		if err != nil {
			if !decode.CheckError(err) {
				return
			}
			logrus.Error(err)
			return
		}

		muxer.Muxer(f, nalus, dts)

		// 将流写入文件
		for _, nalu := range nalus {
			h264file.Write(append([]byte{00, 00, 01}, nalu...))
		}
	})

	if _, err := client.Play(nil); err != nil {
		panic(err)
	}

	client.Wait()
}

func ptr[T any](v T) *T {
	return &v
}

type PSPayloader struct{}

func (p *PSPayloader) Payload(mtu uint16, payload []byte) [][]byte {
	var out [][]byte
	if payload == nil || mtu == 0 {
		return out
	}
	for len(payload) > int(mtu) {
		o := make([]byte, mtu)
		copy(o, payload[:mtu])
		payload = payload[mtu:]
		out = append(out, o)
	}
	o := make([]byte, len(payload))
	copy(o, payload)
	return append(out, o)
}

package ps

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

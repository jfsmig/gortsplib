package formats

import (
	"github.com/pion/rtp"
)

// MPEG2Audio is a format that uses a MPEG-1 or MPEG-2 audio codec.
type MPEG2Audio struct{}

// String implements Format.
func (f *MPEG2Audio) String() string {
	return "MPEG2-audio"
}

// ClockRate implements Format.
func (f *MPEG2Audio) ClockRate() int {
	return 90000
}

// PayloadType implements Format.
func (f *MPEG2Audio) PayloadType() uint8 {
	return 14
}

func (f *MPEG2Audio) unmarshal(
	payloadType uint8, clock string, codec string,
	rtpmap string, fmtp map[string]string,
) error {
	return nil
}

// Marshal implements Format.
func (f *MPEG2Audio) Marshal() (string, map[string]string) {
	return "", nil
}

// PTSEqualsDTS implements Format.
func (f *MPEG2Audio) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}

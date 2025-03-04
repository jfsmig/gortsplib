package formats

import (
	"github.com/pion/rtp"
)

// MPEG2Video is a format that uses a MPEG-1 or MPEG-2 video codec.
type MPEG2Video struct{}

// String implements Format.
func (f *MPEG2Video) String() string {
	return "MPEG2-video"
}

// ClockRate implements Format.
func (f *MPEG2Video) ClockRate() int {
	return 90000
}

// PayloadType implements Format.
func (f *MPEG2Video) PayloadType() uint8 {
	return 32
}

func (f *MPEG2Video) unmarshal(
	payloadType uint8, clock string, codec string,
	rtpmap string, fmtp map[string]string,
) error {
	return nil
}

// Marshal implements Format.
func (f *MPEG2Video) Marshal() (string, map[string]string) {
	return "", nil
}

// PTSEqualsDTS implements Format.
func (f *MPEG2Video) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}

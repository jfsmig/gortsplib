package formats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pion/rtp"
)

func findClockRate(payloadType uint8, rtpMap string) (int, error) {
	// get clock rate from payload type
	// https://en.wikipedia.org/wiki/RTP_payload_formats
	switch payloadType {
	case 0, 1, 2, 3, 4, 5, 7, 8, 9, 12, 13, 15, 18:
		return 8000, nil

	case 6:
		return 16000, nil

	case 10, 11:
		return 44100, nil

	case 14, 25, 26, 28, 31, 32, 33, 34:
		return 90000, nil

	case 16:
		return 11025, nil

	case 17:
		return 22050, nil
	}

	// get clock rate from rtpmap
	// https://tools.ietf.org/html/rfc4566
	// a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
	if rtpMap == "" {
		return 0, fmt.Errorf("attribute 'rtpmap' not found")
	}

	tmp := strings.Split(rtpMap, "/")
	if len(tmp) != 2 && len(tmp) != 3 {
		return 0, fmt.Errorf("invalid rtpmap (%v)", rtpMap)
	}

	v, err := strconv.ParseInt(tmp[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Generic is a generic format.
type Generic struct {
	PayloadTyp uint8
	RTPMap     string
	FMTP       map[string]string

	// clock rate of the format. Filled automatically.
	ClockRat int
}

// Init computes the clock rate of the format. It it mandatory to call it.
func (f *Generic) Init() error {
	f.ClockRat, _ = findClockRate(f.PayloadTyp, f.RTPMap)
	return nil
}

// String returns a description of the format.
func (f *Generic) String() string {
	return "Generic"
}

// ClockRate implements Format.
func (f *Generic) ClockRate() int {
	return f.ClockRat
}

// PayloadType implements Format.
func (f *Generic) PayloadType() uint8 {
	return f.PayloadTyp
}

func (f *Generic) unmarshal(
	payloadType uint8, clock string, codec string,
	rtpmap string, fmtp map[string]string,
) error {
	f.PayloadTyp = payloadType
	f.RTPMap = rtpmap
	f.FMTP = fmtp

	return f.Init()
}

// Marshal implements Format.
func (f *Generic) Marshal() (string, map[string]string) {
	return f.RTPMap, f.FMTP
}

// PTSEqualsDTS implements Format.
func (f *Generic) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}

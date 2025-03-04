package formats

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/pion/rtp"
)

// Vorbis is a format that uses the Vorbis codec.
type Vorbis struct {
	PayloadTyp    uint8
	SampleRate    int
	ChannelCount  int
	Configuration []byte
}

// String implements Format.
func (f *Vorbis) String() string {
	return "Vorbis"
}

// ClockRate implements Format.
func (f *Vorbis) ClockRate() int {
	return f.SampleRate
}

// PayloadType implements Format.
func (f *Vorbis) PayloadType() uint8 {
	return f.PayloadTyp
}

func (f *Vorbis) unmarshal(payloadType uint8, clock string, codec string, rtpmap string, fmtp map[string]string) error {
	f.PayloadTyp = payloadType

	tmp := strings.SplitN(clock, "/", 2)
	if len(tmp) != 2 {
		return fmt.Errorf("invalid clock (%v)", clock)
	}

	sampleRate, err := strconv.ParseInt(tmp[0], 10, 64)
	if err != nil {
		return err
	}
	f.SampleRate = int(sampleRate)

	channelCount, err := strconv.ParseInt(tmp[1], 10, 64)
	if err != nil {
		return err
	}
	f.ChannelCount = int(channelCount)

	for key, val := range fmtp {
		if key == "configuration" {
			conf, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return fmt.Errorf("invalid AAC config (%v)", val)
			}

			f.Configuration = conf
		}
	}

	if f.Configuration == nil {
		return fmt.Errorf("config is missing")
	}

	return nil
}

// Marshal implements Format.
func (f *Vorbis) Marshal() (string, map[string]string) {
	fmtp := map[string]string{
		"configuration": base64.StdEncoding.EncodeToString(f.Configuration),
	}

	return "VORBIS/" + strconv.FormatInt(int64(f.SampleRate), 10) +
			"/" + strconv.FormatInt(int64(f.ChannelCount), 10),
		fmtp
}

// PTSEqualsDTS implements Format.
func (f *Vorbis) PTSEqualsDTS(*rtp.Packet) bool {
	return true
}

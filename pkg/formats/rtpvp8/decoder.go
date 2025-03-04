package rtpvp8

import (
	"errors"
	"fmt"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"

	"github.com/bluenviron/gortsplib/v3/pkg/rtptime"
)

// ErrMorePacketsNeeded is returned when more packets are needed.
var ErrMorePacketsNeeded = errors.New("need more packets")

// ErrNonStartingPacketAndNoPrevious is returned when we received a non-starting
// packet of a fragmented frame and we didn't received anything before.
// It's normal to receive this when we are decoding a stream that has been already
// running for some time.
var ErrNonStartingPacketAndNoPrevious = errors.New(
	"received a non-starting fragment without any previous starting fragment")

// Decoder is a RTP/VP8 decoder.
type Decoder struct {
	timeDecoder         *rtptime.Decoder
	firstPacketReceived bool
	fragments           [][]byte
}

// Init initializes the decoder.
func (d *Decoder) Init() {
	d.timeDecoder = rtptime.NewDecoder(rtpClockRate)
}

// Decode decodes a VP8 frame from a RTP/VP8 packet.
func (d *Decoder) Decode(pkt *rtp.Packet) ([]byte, time.Duration, error) {
	var vpkt codecs.VP8Packet
	_, err := vpkt.Unmarshal(pkt.Payload)
	if err != nil {
		d.fragments = d.fragments[:0] // discard pending fragmented packets
		return nil, 0, err
	}

	if vpkt.PID != 0 {
		d.fragments = d.fragments[:0] // discard pending fragmented packets
		return nil, 0, fmt.Errorf("packets containing single partitions are not supported")
	}

	var frame []byte

	if vpkt.S == 1 {
		d.fragments = d.fragments[:0] // discard pending fragmented packets
		d.firstPacketReceived = true

		if !pkt.Marker {
			d.fragments = append(d.fragments, vpkt.Payload)
			return nil, 0, ErrMorePacketsNeeded
		}

		frame = vpkt.Payload
	} else {
		if len(d.fragments) == 0 {
			if !d.firstPacketReceived {
				return nil, 0, ErrNonStartingPacketAndNoPrevious
			}

			return nil, 0, fmt.Errorf("received a non-starting fragment")
		}

		d.fragments = append(d.fragments, vpkt.Payload)

		if !pkt.Marker {
			return nil, 0, ErrMorePacketsNeeded
		}

		n := 0
		for _, frag := range d.fragments {
			n += len(frag)
		}

		frame = make([]byte, n)
		pos := 0

		for _, frag := range d.fragments {
			pos += copy(frame[pos:], frag)
		}

		d.fragments = d.fragments[:0]
	}

	return frame, d.timeDecoder.Decode(pkt.Timestamp), nil
}

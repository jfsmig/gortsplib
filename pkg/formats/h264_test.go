package formats

import (
	"testing"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

func TestH264Attributes(t *testing.T) {
	format := &H264{
		PayloadTyp:        96,
		SPS:               []byte{0x01, 0x02},
		PPS:               []byte{0x03, 0x04},
		PacketizationMode: 1,
	}
	require.Equal(t, "H264", format.String())
	require.Equal(t, 90000, format.ClockRate())
	require.Equal(t, uint8(96), format.PayloadType())

	sps, pps := format.SafeParams()
	require.Equal(t, []byte{0x01, 0x02}, sps)
	require.Equal(t, []byte{0x03, 0x04}, pps)

	format.SafeSetParams([]byte{0x07, 0x08}, []byte{0x09, 0x0A})

	sps, pps = format.SafeParams()
	require.Equal(t, []byte{0x07, 0x08}, sps)
	require.Equal(t, []byte{0x09, 0x0A}, pps)
}

func TestH264PTSEqualsDTS(t *testing.T) {
	format := &H264{
		PayloadTyp:        96,
		SPS:               []byte{0x01, 0x02},
		PPS:               []byte{0x03, 0x04},
		PacketizationMode: 1,
	}

	require.Equal(t, true, format.PTSEqualsDTS(&rtp.Packet{
		Payload: []byte{0x05},
	}))
	require.Equal(t, false, format.PTSEqualsDTS(&rtp.Packet{
		Payload: []byte{0x01},
	}))
}

func TestH264MediaDescription(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		format := &H264{
			PayloadTyp: 96,
			SPS: []byte{
				0x67, 0x64, 0x00, 0x0c, 0xac, 0x3b, 0x50, 0xb0,
				0x4b, 0x42, 0x00, 0x00, 0x03, 0x00, 0x02, 0x00,
				0x00, 0x03, 0x00, 0x3d, 0x08,
			},
			PPS: []byte{
				0x68, 0xee, 0x3c, 0x80,
			},
			PacketizationMode: 1,
		}

		rtpmap, fmtp := format.Marshal()
		require.Equal(t, "H264/90000", rtpmap)
		require.Equal(t, map[string]string{
			"packetization-mode":   "1",
			"sprop-parameter-sets": "Z2QADKw7ULBLQgAAAwACAAADAD0I,aO48gA==",
			"profile-level-id":     "64000C",
		}, fmtp)
	})

	t.Run("no sps/pps", func(t *testing.T) {
		format := &H264{
			PayloadTyp:        96,
			PacketizationMode: 1,
		}

		rtpmap, fmtp := format.Marshal()
		require.Equal(t, "H264/90000", rtpmap)
		require.Equal(t, map[string]string{
			"packetization-mode": "1",
		}, fmtp)
	})
}

func TestH264DecEncoder(t *testing.T) {
	format := &H264{}

	enc := format.CreateEncoder()
	pkts, err := enc.Encode([][]byte{{0x01, 0x02, 0x03, 0x04}}, 0)
	require.NoError(t, err)
	require.Equal(t, format.PayloadType(), pkts[0].PayloadType)

	dec := format.CreateDecoder()
	byts, _, err := dec.Decode(pkts[0])
	require.NoError(t, err)
	require.Equal(t, [][]byte{{0x01, 0x02, 0x03, 0x04}}, byts)
}

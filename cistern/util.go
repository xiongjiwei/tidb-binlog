package cistern

import (
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/tipb/go-binlog"
)

const lengthOfBinaryTime = 15

// InitLogger initalizes Pump's logger.
func InitLogger(isDebug bool) {
	if isDebug {
		log.SetLevelByString("debug")
	} else {
		log.SetLevelByString("info")
	}
	log.SetHighlighting(false)
}

// ComparePos compares the two positions of binlog items, return 0 when the left equal to the right,
// return -1 if the left is ahead of the right, oppositely return 1.
func ComparePos(left, right binlog.Pos) int {
	if left.Suffix < right.Suffix {
		return -1
	} else if left.Suffix > right.Suffix {
		return 1
	} else if left.Offset < right.Offset {
		return -1
	} else if left.Offset > right.Offset {
		return 1
	} else {
		return 0
	}
}

// CalculateNextPos calculates the position of binlog item next to the given one.
func CalculateNextPos(item binlog.Entity) binlog.Pos {
	pos := item.Pos
	// 4 bytes(magic) + 8 bytes(size) + length of payload + 4 bytes(CRC)
	pos.Offset += int64(len(item.Payload) + 16)
	return pos
}

func encodePayload(payload []byte) ([]byte, error) {
	nowBinary, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, errors.Trace(err)
	}
	n1 := len(nowBinary)
	n2 := len(payload)
	data := make([]byte, n1+n2)
	copy(data[:n1], nowBinary)
	copy(data[n1:], payload)
	return data, nil
}

func decodePayload(value []byte) ([]byte, time.Duration, error) {
	var ts time.Time
	n1 := lengthOfBinaryTime

	data := make([]byte, len(value))
	copy(data, value)

	if err := ts.UnmarshalBinary(data[:n1]); err != nil {
		return nil, 0, errors.Trace(err)
	}
	return data[n1:], time.Now().Sub(ts), nil
}

// combine suffix offset in one float, the format would be suffix.offset
func posToFloat(pos *binlog.Pos) float64 {
	var decimal float64
	decimal = float64(pos.Offset)
	for decimal > 1 {
		decimal = decimal / 10
	}
	return float64(pos.Suffix) + decimal
}
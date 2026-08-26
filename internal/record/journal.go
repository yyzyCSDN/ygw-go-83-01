package record

import (
	"encoding/binary"
)

func encodeEvent(e Event) []byte {
	payload, err := marshalEvent(e)
	if err != nil {
		return []byte{'\n'}
	}
	sum := digest(payload)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, sum)
	line := make([]byte, 0, len(payload)+8+1)
	line = append(line, buf...)
	line = append(line, payload...)
	line = append(line, '\n')
	return line
}

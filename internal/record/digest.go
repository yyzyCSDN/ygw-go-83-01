package record

import "github.com/cespare/xxhash/v2"

func digest(payload []byte) uint64 {
	return xxhash.Sum64(payload)
}

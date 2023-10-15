package conditions

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"strconv"
)

func Random() int64 {
	// Seeding not needed as of go 1.20
	return rand.Int63()
}

var SessionRand = rand.Int63()

func SessionRandom() int64 {
	return SessionRand
}

func RandomForKey(key string, seed int64) int64 {
	hasher := sha256.New()
	hasher.Write([]byte(key + strconv.FormatInt(seed, 16)))
	sum := hasher.Sum(nil)
	uint := binary.LittleEndian.Uint64(sum[0:8])
	// 63 bit positive int like rand.Int63()
	return int64(uint % (1 << 62))
}

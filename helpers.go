package main

import (
	"math/rand"
	"strings"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	idxBits = 6              // 6 bits to represent a letter index
	mask    = 1<<idxBits - 1 // All 1-bits, as many as idxBits
	idxMax  = 63 / idxBits   // # of letter indices fitting in 63 bits
)

// GetRandomFileName generates a random alphanumeric filename with length Config.Size,
// ending in the passed in extention.
func GetRandomFileName(extension string) (fileName string, err error) {
	sb := strings.Builder{}
	sb.Grow(Config.Size)

	for i, cache, remain := Config.Size-1, src.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxMax
		}

		if idx := int(cache & mask); idx < len(alphabet) {
			err = sb.WriteByte(alphabet[idx])
			if err != nil {
				return fileName, err
			}
			i--
		}

		cache >>= idxBits
		remain--
	}

	_, err = sb.WriteString("." + extension)
	if err != nil {
		return fileName, err
	}

	fileName = sb.String()

	return fileName, err
}

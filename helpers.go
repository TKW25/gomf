package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
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
	sb.Grow(Config.FileNameLength)

	for i, cache, remain := Config.FileNameLength-1, src.Int63(), idxMax; i >= 0; {
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

	_, err = sb.WriteString(extension)
	if err != nil {
		return fileName, err
	}

	fileName = sb.String()

	return fileName, err
}

func GetHash(file *multipart.File) (string, error) {
	hash := sha1.New()
	if _, err := io.Copy(hash, *file); err != nil {
		return "", err
	}

	if _, err := (*file).Seek(0, 0); err != nil {
		log.Println(err)
		log.Printf("Failed seeking to the start of the file")
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

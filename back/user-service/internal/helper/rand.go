package helper

import (
	"crypto/sha256"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

func GeneralDuration(times time.Duration, minTime, maxTime int, duration time.Duration) time.Duration {
	return times + time.Duration(+RandInt(minTime, maxTime))*duration
}

func GetRandString() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	h := sha256.New()
	h.Write(b)

	return string(h.Sum(nil))
}

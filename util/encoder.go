package util

import (
	"math/rand"
)

func GenerateShortenURL() string {
    const dict = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
    const length = 8

    shortURL := make([]byte, length)
    for i := range shortURL {
        shortURL[i] = dict[rand.Intn(len(dict))]
    }
    return string(shortURL)
}

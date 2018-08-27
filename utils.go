package main

import (
	"crypto/rand"
	"fmt"
)

// GenereteID generating unique id
func GenereteID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

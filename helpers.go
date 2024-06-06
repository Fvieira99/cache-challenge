package main

import "math/rand"

func GenerateRandomID(i int) int {
	if i < 100 {
		return i + 1
	}
	return rand.Intn(99) + 1
}

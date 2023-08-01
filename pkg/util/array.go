package util

import (
	"math/rand"
	"time"
)

func ShuffleArray(arr []interface{}) {
	rand.Seed(time.Now().UnixNano())
	n := len(arr)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

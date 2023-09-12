package utils

import (
"fmt"
"math/rand"
)

func GetRandomInt(min, max int) int {
	return rand.Intn(max - min) + min
}

func GetRandomStringNumber(len int) string {
	number := ""
	for i := 0; i < len; i++ {
		number = fmt.Sprint(number, GetRandomInt(0, 9))
	}
	return number
}

package test

import (
	"log"
	"testing"
)

func TestDir(t *testing.T) {

	a := []int{1, 2, 3, 4, 5}

	for i, i2 := range a {
		log.Println(i, i2)
	}

}

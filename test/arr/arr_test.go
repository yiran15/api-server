package arr_test

import (
	"fmt"
	"testing"
)

func TestArr(t *testing.T) {
	var a []int
	b := []int{}

	fmt.Println(a == nil)
	fmt.Println(b == nil)

	fmt.Println(len(a))
	fmt.Println(len(b))
}

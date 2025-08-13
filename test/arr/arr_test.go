package arr_test

import (
	"fmt"
	"strings"
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

func TestSplit(t *testing.T) {
	a := "/api/v1/user/login"
	a = strings.TrimPrefix(a, "/")
	ty := strings.Split(a, "/")[2]
	fmt.Println(ty)
}

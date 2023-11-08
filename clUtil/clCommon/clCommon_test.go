package clCommon

import (
	"fmt"
	"testing"
)

func TestUnderlineToUppercase(t *testing.T) {

	fmt.Printf("%v\n", UnderlineToUppercase(true, "hello_world"))
	fmt.Printf("%v\n", UnderlineToUppercase(false, "hello_world"))

}
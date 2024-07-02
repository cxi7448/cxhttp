package clCommon

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnderlineToUppercase(t *testing.T) {
	str := ";base64,"
	index := strings.Index(str, ";base64,")
	fmt.Println(str[index:])
}

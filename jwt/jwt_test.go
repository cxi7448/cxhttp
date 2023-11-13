package jwt

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(nil)
	fmt.Println(err)
	fmt.Println(token)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySW5mbyI6bnVsbCwiZXhwIjoxNjk1MTExNDQ2LCJpc3MiOiJnb2FwaSIsInN1YiI6InRlc3QiLCJDcmVhdGVUaW1lIjoxNjk1MTExNDQ2LCJSZWZsdXNoVGltZSI6MTY5NTExMTQ0Nn0.aBU50-dRXy4LCJy_Kj_YHsgz9bsmysnFOkh-ZIdg_uY
}

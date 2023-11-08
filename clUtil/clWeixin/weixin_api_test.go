package clWeixin

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"testing"
)

const APP_ID = "wxe28b0753361967f4"
const APP_SECRET = "96472450cf326f355e34509c9865de29"

func TestCode2Session(t *testing.T) {
	obj, _ := Code2Session(APP_ID, APP_SECRET, "0939pPkl2QxUv84eAbol2m2IhQ09pPkI")
	clLog.Info("%+v", obj)
}

func TestGetAccessToken(t *testing.T) {
	obj, err := GetAccessToken(APP_ID, APP_SECRET)
	clLog.Info("%v, %v", obj, err)
}

func TestDecryptUserInfo(t *testing.T) {
	data := `W93DvXcLpGwiHmtnukZ0V6u/yu9WU4Ag2XyclaexZXdkJQT21toPibUpFm9 d/9tUo718XC0i6OUbKBkb88uroOvtPln7qO6Ku8d5CSttQgERovrWz7LjucTKCHny UIC0NlODKzO7D1ICrPBj0Qd6IqrjAdn79MzDmxLJvSWpWGzWkoHkcGX00Vm60pxN2RofroEpfn5eKs mE H7mdzqQ2UX6lZmlbLt6n 48OtYPJg2vk0o4KOF3p6TuUuJgXqHO5/6OO/whqL/8/bB7l53dUpDcE7bENYCxoSKabjM56EJTaVB7MTrqw8WHP9UbANHKw86Fz/jeV9W6fxwy0fpitBHUoTvxFuB0kUzk6lN/P2eyThx7gLDhtkWPgM5b1qgOQwDIMCmWqv/9RNghXw8nvTOmmkUo2DsU50ZQSOJk=`
	key := `stnYVNdp1 kcAXZgMCOlWw==`
	iv := `Xo1iejT8EpXvFxkkQxUTPw==`

	DecryptUserInfo(data, key, iv)
}

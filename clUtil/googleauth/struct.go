package googleauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func New() *GoogleAuth {
	return &GoogleAuth{}
}

func (this *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (this *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (this *GoogleAuth) base32Encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (this *GoogleAuth) base32Decode(src string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(src)
}
func (this *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (this *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) + (uint32(bts[2]) << 8) + uint32(bts[3])
}

func (this *GoogleAuth) oneTimePassword(key, data []byte) uint32 {
	hash := this.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := this.toUint32(hashParts)
	return number % 1000000
}

func (this *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, this.un())
	return strings.ToUpper(this.base32Encode(this.hmacSha1(buf.Bytes(), nil)))
}

func (this *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := this.base32Decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := this.oneTimePassword(secretKey, this.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}
func (this *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

func (this *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	return ""
}
func (this *GoogleAuth) Verify(secret, code string) bool {
	_code, err := this.GetCode(secret)
	if err != nil {
		clLog.Error("谷歌身份验证器验证失败:%v", err)
		return false
	}
	return _code == code
}

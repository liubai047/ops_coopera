package crand

import (
	"math/rand"
	"strings"
	"time"
)

type mRand struct {
	charTyp charType
}

func NewMRand(charTyp charType) *mRand {
	return &mRand{
		charTyp: charTyp,
	}
}

const numberChar = "123465789"
const lowAZ = "abcdefghijklmnopqrstuvwxyz"
const upAZ = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func stringWithCharset(length int, charset string) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type charType int

const (
	LowChar    charType = iota // 小写字符集
	UpChar     charType = iota // 大写字符集
	NumberChar charType = iota // 数字集
	UpLowChar  charType = iota // 大小写混合字符集
	AllChar    charType = iota // 大小写字符+数字字符集
)

func (r *mRand) RandomString(length int) string {
	var charset strings.Builder
	switch r.charTyp {
	case LowChar:
		charset.WriteString(lowAZ)
	case UpChar:
		charset.WriteString(upAZ)
	case NumberChar:
		charset.WriteString(numberChar)
	case UpLowChar:
		charset.WriteString(lowAZ)
		charset.WriteString(upAZ)
	case AllChar:
		charset.WriteString(lowAZ)
		charset.WriteString(upAZ)
		charset.WriteString(numberChar)
	}
	return stringWithCharset(length, charset.String())
}

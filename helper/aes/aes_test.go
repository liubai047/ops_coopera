package aes

import (
	"testing"
)

func Test_AesBase64Encrypt(t *testing.T) {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	var base64Text = "DQIuAJkBADOJ8v2+Abs71Nfut473cteCNkoRgN3S0rE="
	res, err := Base64Encrypt(text, aesKey)
	if res != base64Text || err != nil {
		t.Errorf("base64加密方法和预期值不一致")
	}
}

func Test_AesBase64Decrypt(t *testing.T) {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	var base64Text = "DQIuAJkBADOJ8v2+Abs71Nfut473cteCNkoRgN3S0rE="
	res, err := Base64Decrypt(base64Text, aesKey)
	if string(res) != text || err != nil {
		t.Errorf("base64解密方法和预期值不一致")
	}
}

func Test_AesHexEncrypt(t *testing.T) {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	var hexText = "0d022e009901003389f2fdbe01bb3bd4d7eeb78ef772d782364a1180ddd2d2b1"
	res, err := HexEncrypt(text, aesKey)
	if string(res) != hexText || err != nil {
		t.Errorf("hex加密方法和预期值不一致")
	}
}

func Test_AesHexDecrypt(t *testing.T) {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	var hexText = "0d022e009901003389f2fdbe01bb3bd4d7eeb78ef772d782364a1180ddd2d2b1"
	res, err := HexDecrypt(hexText, aesKey)
	if string(res) != text || err != nil {
		t.Errorf("hex解密方法和预期值不一致")
	}
}

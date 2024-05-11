package aes

import "fmt"

func ExampleBase64Encrypt() {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	fmt.Println(Base64Encrypt(text, aesKey))
	// Output:
	// DQIuAJkBADOJ8v2+Abs71Nfut473cteCNkoRgN3S0rE=, nil
}

func ExampleBase64Decrypt() {
	var aesKey = "sdfnisndfESVSsxx"
	var base64Text = "DQIuAJkBADOJ8v2+Abs71Nfut473cteCNkoRgN3S0rE="
	fmt.Println(Base64Decrypt(base64Text, aesKey))
	// Output:
	// 这是一段文本哈, nil
}

func ExampleHexEncrypt() {
	var aesKey = "sdfnisndfESVSsxx"
	var text = "这是一段文本哈"
	fmt.Println(HexEncrypt(text, aesKey))
	// Output:
	// 0d022e009901003389f2fdbe01bb3bd4d7eeb78ef772d782364a1180ddd2d2b1, nil
}

func ExampleHexDecrypt() {
	var aesKey = "sdfnisndfESVSsxx"
	var hexText = "0d022e009901003389f2fdbe01bb3bd4d7eeb78ef772d782364a1180ddd2d2b1"
	fmt.Println(HexDecrypt(hexText, aesKey))
	// Output:
	// 这是一段文本哈, nil
}

package server

import (
	"strings"
)

const (
	prefBarCode = "BtnProperties/BarCode/"
	prefWheight = "BtnProperties/BtnWheight/"
	prefVolume  = "BtnProperties/BtnVolume/"
	prefBtn     = "Btn"
)

func getBtnBloc(name string) string {

	if strings.Contains(name, "Wheight") || strings.Contains(name, "Volume") || strings.Contains(name, "Btn796") || strings.Contains(name, "Btn006") || strings.Contains(name, "Btn778") || strings.Contains(name, "Btn123") {
		return "BtnProperties"
	}

	if strings.Contains(name, "Btn1") || strings.Contains(name, "Btn2") || strings.Contains(name, "Btn1") || strings.Contains(name, "Btn3") || strings.Contains(name, "Btn4") ||
		strings.Contains(name, "Btn5") || strings.Contains(name, "Btn6") || strings.Contains(name, "Btn7") || strings.Contains(name, "Btn8") ||
		strings.Contains(name, "Btn9") || strings.Contains(name, "Btn10") || strings.Contains(name, "Btn11") || strings.Contains(name, "Btn12") ||
		strings.Contains(name, "Btn13") || strings.Contains(name, "Btn14") || strings.Contains(name, "Btn15") || strings.Contains(name, "Btn16") {
		return "BtnOrder"
	}

	return ""

}

func setBtnOrders() map[string]string {
	menuOrder := make(map[string]string, 16)
	menuOrder["Btn11"] = "АЭ"
	menuOrder["Btn12"] = "Чек"
	menuOrder["Btn14"] = "РНк"
	menuOrder["Btn21"] = "РНк/у"
	menuOrder["Btn22"] = "РНр/у"
	menuOrder["Btn23"] = "РнБ"
	menuOrder["Btn24"] = "Прм/у"
	menuOrder["Btn31"] = "ИМ1"
	menuOrder["Btn32"] = "ИМ2"
	menuOrder["Btn33"] = "ИМ3"
	menuOrder["Btn34"] = "ИМ4"
	menuOrder["Btn41"] = "ИП1"
	menuOrder["Btn42"] = "ИП2"
	menuOrder["Btn43"] = "ИП3"
	menuOrder["Btn44"] = "ИП4"

	return menuOrder
}

func textForUser(text string) string {
	switch {
	case strings.Contains(text, "BarCode"):
		return "ШК"
	case strings.Contains(text, "Wheight"):
		return "вес"
	case strings.Contains(text, "Volume"):
		return "объём"
	}

	return ""
}

func getCode(text string) string {
	code := text
	prefixies := []string{prefBarCode, prefWheight, prefVolume, prefBtn}

	for _, v := range prefixies {
		code = strings.Replace(code, v, "", -1)
	}
	return code
}

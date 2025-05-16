package debug

import (
	"fmt"
	"github.com/fatih/color"
)

func PrintError(err error) {
	errLine := fmt.Sprintf("[DEBUG-ERROR] %s", err)
	fmt.Println(color.RedString(errLine))
}

func PrintJson(json string) {
	jsonLine := fmt.Sprintf("[DEBUG-JSON] %s", json)
	fmt.Print(color.GreenString(jsonLine))
}

func PrintInfo(info string) {
	infoLine := fmt.Sprintf("[DEBUG-INFO] %s", info)
	fmt.Println(color.BlueString(infoLine))
}

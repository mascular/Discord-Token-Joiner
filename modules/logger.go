package headers

import (
	"fmt"

	"github.com/fatih/color"
)

func MaskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:12] + "..." + token[len(token)-12:]
}

func LogSuccess(token string, statusCode int) {
	green := color.New(color.FgGreen).SprintFunc()
	masked := MaskToken(token)
	fmt.Printf("%s %s | Status Code: %d\n", green("[+]"), masked, statusCode)
}

func LogFailure(token string, statusCode int) {
	red := color.New(color.FgRed).SprintFunc()
	masked := MaskToken(token)
	fmt.Printf("%s %s | Status Code: %d\n", red("[-]"), masked, statusCode)
}

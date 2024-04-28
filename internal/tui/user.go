package tui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func YesNoLoop(question string) bool {
	yesRegexp, _ := regexp.Compile("[yY][eE]?[sS]?")
	noRegexp, _ := regexp.Compile("[nN][oO]?")

	for {
		fmt.Printf("%v (y/n):", question)

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r\n", "", -1)

		matchYes := yesRegexp.MatchString(text)
		matchNo := noRegexp.MatchString(text)

		if matchYes {
			return true

		} else if matchNo {
			return false

		}
	}
}

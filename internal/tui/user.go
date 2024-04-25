package tui

import (
	"os"
	"strings"
	"bufio"
	"fmt"
	"regexp"
)

func YesNoLoop(question string) bool {
	for true {
		fmt.Printf("%v (y/n):", question)

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r\n", "", -1)

		matchYes, _ := regexp.MatchString("[yY][eE]?[sS]?", text)
		matchNo, _ := regexp.MatchString("[nN][oO]?", text) 

		if matchYes {
			return true;

		} else if matchNo {
			return false;

		}
	}

	return true
}

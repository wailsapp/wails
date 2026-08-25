package exec

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	return line == "y" || line == "Y" || line == "yes"
}

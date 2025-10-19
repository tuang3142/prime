package termbell

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/term"
)

func Run() {
	fd := int(syscall.Stdin)
	if !term.IsTerminal(fd) {
		fmt.Println("stdin is not a terminal")
		os.Exit(1)
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("failed to enable raw mode:", err)
		os.Exit(1)
	}
	defer term.Restore(fd, state)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a digit (0–9) to beep that many times. Ctrl+C to quit.")

	for {
		c, err := reader.ReadByte()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("read error:", err)
			break
		}

		if c >= '0' && c <= '9' {
			n := int(c - '0')
			for i := 0; i < n; i++ {
				fmt.Print("beep\a\r\n")
				time.Sleep(time.Second)
			}
		}
	}
}

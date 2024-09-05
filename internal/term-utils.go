package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type Position struct {
	X int
	Y int
}

const esc = '\x1B'

func CliCommandSave() string {
	return "\u001b7"
}

func CliCommandRestore() string {
	return "\u001b8"
}

func CliCommandStartOfLine() string {
	return "\u001b8"

}
func CliCommandUp(count int) string {
	if count > 0 {
		return fmt.Sprintf("%s[%dA", string(esc), count)
	}
	return ""
}

func GetCursorPos() (*Position, error) {
	reader := bufio.NewReader(os.Stdin)
	res := make(chan string, 1)
	go func() {
		termState, _ := term.GetState(int(os.Stdin.Fd()))
		term.MakeRaw(int(os.Stdin.Fd()))

		reader.ReadBytes('\x1B')
		read, _ := reader.ReadString('R')

		term.Restore(int(os.Stdin.Fd()), termState)
		res <- read
	}()
	go os.Stdout.WriteString("\x1B[6n")

	result := <-res
	regex, err := regexp.Compile(`\[(\d+);(\d+)R`)
	match := regex.FindStringSubmatch(result)
	if len(match) > 1 {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		return &Position{x, y}, nil
	}
	return nil, err
}

var osStdin = os.Stdin

func EnableStdinEcho() error {
	const ioctlReadTermios = unix.TIOCGETA
	const ioctlWriteTermios = unix.TIOCSETA

	termios, err := unix.IoctlGetTermios(int(osStdin.Fd()), ioctlReadTermios)
	if err != nil {
		return err
	}

	newState := *termios
	newState.Lflag |= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(int(osStdin.Fd()), ioctlWriteTermios, &newState); err != nil {
		return err
	}

	return nil
}

func DisableStdinEcho() error {
	const ioctlReadTermios = unix.TIOCGETA
	const ioctlWriteTermios = unix.TIOCSETA

	termios, err := unix.IoctlGetTermios(int(osStdin.Fd()), ioctlReadTermios)
	if err != nil {
		fmt.Println("Oops", err)
		return err
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(int(osStdin.Fd()), ioctlWriteTermios, &newState); err != nil {
		fmt.Println("Oops2", err)
		return err
	}

	return nil
}

func StdinEchoEnabled() bool {
	const ioctlReadTermios = unix.TIOCGETA

	termios, _ := unix.IoctlGetTermios(int(osStdin.Fd()), ioctlReadTermios)
	return termios.Lflag&unix.ECHO == unix.ECHO
}

package node

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

type Screen struct {
	Width  int
	Height int
}

func (s *Screen) Clear() {
	fmt.Print("\033[H\033[2J") // clear screen
}

func (s *Screen) Write(msg string, a ...any) {
	fmt.Print("\033[2K\r")
	fmt.Printf(msg, a...)
	fmt.Print("[YOU]>")
}

var Console Screen

func init() {
	width, height, _ := terminal.GetSize(0)
	Console = Screen{Width: width, Height: height}

}

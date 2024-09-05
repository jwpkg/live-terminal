package main

import (
	living_terminal "github.com/living-terminal"
	"golang.org/x/term"
)

type SinusWaveComponent struct {
	Text            string
	reRenderChannel chan bool
	CurrentFrame    int
}

func (sinus *SinusWaveComponent) Init(requestReRender chan bool) {
	sinus.reRenderChannel = requestReRender
}

func (sinus *SinusWaveComponent) Render() string {
	runes := []rune(sinus.Text)

	size, _, _ := term.GetSize(int(living_terminal.OriginalStdout.Fd()))

	result := make([]rune, 0)
	for pos := range size {
		result = append(result, runes[(pos+sinus.CurrentFrame)%len(runes)])
	}
	return string(result)
}

func (scroller *SinusWaveComponent) Finish() {
	scroller.reRenderChannel = nil
}

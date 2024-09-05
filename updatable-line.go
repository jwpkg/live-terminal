package living_terminal

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/living-terminal/internal"
	"golang.org/x/term"
)

type renderLineUpdate struct {
	startLinePosition int
	text              string
}

type sharedLineStruct struct {
	lineCounter          *internal.LineCounter
	usageCount           int
	renderLineUpdateChan chan renderLineUpdate
	done                 chan bool
	mu                   sync.Mutex
}

var sharedLine = sharedLineStruct{}

func retain() {
	sharedLine.mu.Lock()
	defer sharedLine.mu.Unlock()

	if sharedLine.usageCount == 0 {
		sharedLine.lineCounter = internal.NewLineCounter()
		sharedLine.renderLineUpdateChan = make(chan renderLineUpdate, 200)
		sharedLine.done = make(chan bool)
		go renderer()
	}

	sharedLine.usageCount++
}

func release() {
	sharedLine.mu.Lock()
	defer sharedLine.mu.Unlock()

	if sharedLine.usageCount > 0 {
		sharedLine.usageCount--

		if sharedLine.usageCount == 0 {
			close(sharedLine.renderLineUpdateChan)
			<-sharedLine.done

			sharedLine.lineCounter.Stop()
			sharedLine.renderLineUpdateChan = nil
			sharedLine.lineCounter = nil
			sharedLine.done = nil
		}
	}
}

func renderer() {
	render := func(renderUpdate renderLineUpdate) {
		_, height, _ := term.GetSize(0)

		offSet := sharedLine.lineCounter.Count() - renderUpdate.startLinePosition

		if offSet < height { // check off screen
			outputString := internal.CliCommandSave() +
				internal.CliCommandUp(offSet) + "\r" +
				renderUpdate.text +
				internal.CliCommandRestore()

			os.Stdout.WriteString(outputString)
		}
	}

	for {
		renderUpdate, ok := <-sharedLine.renderLineUpdateChan
		if ok {
			render(renderUpdate)
		} else {
			sharedLine.done <- true
			return
		}
	}
}

// Updatable line in the terminal
type Line struct {
	startPos    int
	currentText string
}

func NewLine(startText string) *Line {
	retain()

	fmt.Println(startText)

	line := &Line{
		startPos:    sharedLine.lineCounter.Count() - 1,
		currentText: startText,
	}

	return line
}

func (lineUpdater *Line) Finish() {

	release()
}

func (lineUpdater *Line) Update(newText string) {
	if newText != lineUpdater.currentText {
		paddingLength := max(0, len(lineUpdater.currentText)-len(newText))
		lineUpdater.currentText = newText

		renderUpdate := renderLineUpdate{
			startLinePosition: lineUpdater.startPos,
			text:              newText + strings.Repeat(" ", paddingLength),
		}

		sharedLine.renderLineUpdateChan <- renderUpdate
	}
}

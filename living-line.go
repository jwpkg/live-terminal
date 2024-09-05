package living_terminal

import (
	"os"
	"time"
)

var OriginalStdin = os.Stdin
var OriginalStdout = os.Stdout
var OriginalStderr = os.Stderr

// Updatable line in the terminal
type LivingLine struct {
	reRenderChannel chan bool
	components      []LivingComponent
	line            Line
	done            chan bool
}

const renderInterval = time.Millisecond * 10

func NewLivingLine(components ...LivingComponent) *LivingLine {
	result := LivingLine{
		reRenderChannel: make(chan bool, 1),
		done:            make(chan bool),
		components:      components,
		line:            *NewLine(renderComponents(components)),
	}

	for _, component := range components {
		component.Init(result.reRenderChannel)
	}

	go result.run()

	return &result
}

func (line *LivingLine) Update(components ...LivingComponent) {
	for _, component := range line.components {
		component.Finish()
	}
	line.components = components

	line.line.Update(renderComponents(components))

	for _, component := range components {
		component.Init(line.reRenderChannel)
	}
}

func renderComponents(components []LivingComponent) string {
	result := ""

	for _, component := range components {
		result += component.Render()
	}
	return result
}

func (line *LivingLine) run() {
	ticker := time.NewTicker(renderInterval)
	needsReRender := false
	ok := true
	for {
		select {
		case needsReRender, ok = <-line.reRenderChannel:
			if !ok {
				ticker.Stop()
				line.done <- true
				return
			}
		case <-ticker.C:
			if needsReRender {
				needsReRender = false
				line.line.Update(renderComponents(line.components))
			}
		}
	}
}

func (line *LivingLine) Finish() {
	for _, component := range line.components {
		component.Finish()
	}

	close(line.reRenderChannel)
	<-line.done
	line.line.Finish()
}

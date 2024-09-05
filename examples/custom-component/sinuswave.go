package main

import (
	"math"
	"time"

	living_terminal "github.com/jwpkg/living-terminal"
)

type SinusWave struct {
	lines        []living_terminal.LivingLine
	components   []*SinusWaveComponent
	currentFrame int
	finished     chan bool
	done         chan bool
}

func NewSinusWave() *SinusWave {
	lines := make([]string, 0)

	for y := 1.0; y >= -1.0; y -= 0.2 {
		currentLine := ""
		for x := 0.0; x < math.Pi*2-0.1; x += 0.1 {
			sin := math.Sin(x)
			if (y+0.1) > sin && (y-0.1) <= sin {
				currentLine += "*"
			} else {
				currentLine += " "
			}
		}
		lines = append(lines, currentLine)
	}

	components := make([]*SinusWaveComponent, 0)
	livingLines := make([]living_terminal.LivingLine, 0)
	for _, line := range lines {
		component := SinusWaveComponent{
			Text: line,
		}
		livingLine := living_terminal.NewLivingLine(&component)

		components = append(components, &component)
		livingLines = append(livingLines, *livingLine)
	}

	sinusWave := &SinusWave{
		components: components,
		lines:      livingLines,
		finished:   make(chan bool),
		done:       make(chan bool),
	}

	go sinusWave.run()

	return sinusWave
}

func (sinusWave *SinusWave) run() {
	ticker := time.NewTicker(time.Millisecond * 30)

	for {
		select {
		case <-sinusWave.finished:
			ticker.Stop()
			sinusWave.done <- true
			return
		case <-ticker.C:
			sinusWave.currentFrame++
			for _, component := range sinusWave.components {
				component.CurrentFrame = sinusWave.currentFrame
				component.reRenderChannel <- true
			}
		}
	}
}

func (sinusWave *SinusWave) Finish() {
	sinusWave.finished <- true
	<-sinusWave.done
	for _, line := range sinusWave.lines {
		line.Finish()
	}
}

package main

import (
	"fmt"
	"time"

	spinner "github.com/gabe565/go-spinners"
	living_terminal "github.com/living-terminal"
	components "github.com/living-terminal/components"
)

func main() {
	fmt.Println("Start")

	progressBar := components.NewProgressBar()
	line1 := living_terminal.NewLivingLine(
		progressBar,
	)
	defer line1.Finish()

	line2 := living_terminal.NewLivingLine(
		components.NewLivingSpinner(spinner.Dots),
		components.NewLivingText(" Processing"),
		components.NewLivingSpinner(spinner.SimpleDots),
	)
	defer line2.Finish()

	line3 := living_terminal.NewLivingLine(
		&components.LivingScroller{
			Text:     "Scrolling text! - ",
			Interval: time.Millisecond * 50,
			Size:     50,
		},
	)
	defer line3.Finish()

	fmt.Println("Just print below")

	for i := range 100 {
		progressBar.Update(i + 1)
		time.Sleep(time.Millisecond * 50)
	}

	line2.Update(components.NewLivingText("Processing Done!"))
}

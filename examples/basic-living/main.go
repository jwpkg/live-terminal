package main

import (
	"fmt"
	"time"

	spinner "github.com/gabe565/go-spinners"
	living_terminal "github.com/jwpkg/living-terminal"
)

func main() {
	fmt.Println("Start")

	progressBar := living_terminal.NewProgressBar()
	line1 := living_terminal.NewLivingLine(
		progressBar,
	)
	defer line1.Finish()

	line2 := living_terminal.NewLivingLine(
		living_terminal.NewLivingSpinner(spinner.Dots),
		living_terminal.NewLivingText(" Processing"),
		living_terminal.NewLivingSpinner(spinner.SimpleDots),
	)
	defer line2.Finish()

	line3 := living_terminal.NewLivingLine(
		&living_terminal.LivingScroller{
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

	line2.Update(living_terminal.NewLivingText("Processing Done!"))
}

package components

import (
	"time"

	spinner "github.com/gabe565/go-spinners"
)

type LivingSpinner struct {
	Spinner         spinner.Spinner
	reRenderChannel chan bool
	currentFrame    int
	done            chan bool
	finished        chan bool
}

func NewLivingSpinner(spinner spinner.Spinner) *LivingSpinner {
	return &LivingSpinner{
		Spinner: spinner,
	}
}

func (spinner *LivingSpinner) Init(requestReRender chan bool) {
	spinner.reRenderChannel = requestReRender
	spinner.done = make(chan bool)
	spinner.finished = make(chan bool)
	go spinner.run()
}

func (spinner *LivingSpinner) run() {
	ticker := time.NewTicker(spinner.Spinner.Interval)
	frameCount := len(spinner.Spinner.Frames)

	for {
		select {
		case <-spinner.finished:
			ticker.Stop()
			spinner.done <- true
			return
		case <-ticker.C:
			spinner.currentFrame = (spinner.currentFrame + 1) % frameCount
			spinner.reRenderChannel <- true
		}
	}
}

func (spinner *LivingSpinner) Render() string {
	return spinner.Spinner.Frames[spinner.currentFrame]
}

func (spinner *LivingSpinner) Finish() {
	spinner.finished <- true
	<-spinner.done
	spinner.done = nil
	spinner.reRenderChannel = nil
}

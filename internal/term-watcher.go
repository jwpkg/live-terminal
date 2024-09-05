package internal

import (
	"os"
	"time"

	"golang.org/x/term"
)

// T
type TermWatcher struct {
	// Channel to recieve width updates
	WidthChan chan int

	// Channel to recieve height updates
	HeightChan chan int

	// Channel to recieve echo toggle updates
	StdinEchoChan chan bool

	// Current terminal Width
	Width int

	// Current terminal Height
	Height int

	// Current terminal stdin echo
	StdinEcho bool

	// Notify routine to stop
	done chan bool

	// Current running state
	active bool
}

const pollingInterval = time.Millisecond * 100

func NewTermWatcher() *TermWatcher {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		panic(err)
	}

	return &TermWatcher{
		WidthChan:     make(chan int, 1),
		HeightChan:    make(chan int, 1),
		StdinEchoChan: make(chan bool, 1),
		done:          make(chan bool, 1),
		active:        false,

		Width:     width,
		Height:    height,
		StdinEcho: StdinEchoEnabled(),
	}
}

func (termOptions *TermWatcher) Start() {
	if termOptions.active {
		return
	}
	go termOptions.run()
}

func (termOptions *TermWatcher) Stop() {
	if !termOptions.active {
		return
	}

	termOptions.done <- true
	termOptions.active = false
}

func (termOptions *TermWatcher) run() {
	ticker := time.NewTicker(pollingInterval)
	for {
		select {
		case <-termOptions.done:
			ticker.Stop()
			return
		case <-ticker.C:

			width, height, _ := term.GetSize(int(os.Stdout.Fd()))
			if width != termOptions.Width {
				termOptions.Width = width
				termOptions.WidthChan <- width
			}

			if height != termOptions.Height {
				termOptions.Height = height
				termOptions.HeightChan <- height
			}

			echoOn := StdinEchoEnabled()
			if echoOn != termOptions.StdinEcho {
				termOptions.StdinEcho = echoOn
				termOptions.StdinEchoChan <- echoOn
			}
		}
	}
}

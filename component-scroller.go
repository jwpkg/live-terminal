package living_terminal

import (
	"time"

	"golang.org/x/term"
)

type LivingScroller struct {
	Text            string
	Size            int
	reRenderChannel chan bool
	CurrentFrame    int
	Interval        time.Duration
	done            chan bool
	finished        chan bool
}

func NewLivingScroller(text string) *LivingScroller {
	return &LivingScroller{
		Text:     text,
		Interval: time.Millisecond * 50,
	}
}

func (scroller *LivingScroller) Init(requestReRender chan bool) {
	scroller.reRenderChannel = requestReRender
	scroller.done = make(chan bool)
	scroller.finished = make(chan bool)
	go scroller.run()
}

func (scroller *LivingScroller) run() {
	ticker := time.NewTicker(scroller.Interval)
	frameCount := len(scroller.Text)

	for {
		select {
		case <-scroller.finished:
			ticker.Stop()
			scroller.done <- true
			return
		case <-ticker.C:
			scroller.CurrentFrame = (scroller.CurrentFrame + 1) % frameCount
			scroller.reRenderChannel <- true
		}
	}
}

func (scroller *LivingScroller) Render() string {
	runes := []rune(scroller.Text)
	size := scroller.Size
	if size == 0 {
		size, _, _ = term.GetSize(int(OriginalStdout.Fd()))
	}

	result := make([]rune, 0)
	for pos := range size {
		result = append(result, runes[(pos+scroller.CurrentFrame)%len(runes)])
	}
	return string(result)
}

func (scroller *LivingScroller) Finish() {
	scroller.finished <- true
	<-scroller.done
	scroller.done = nil
	scroller.reRenderChannel = nil
}

package components

import (
	"fmt"
	"math"
	"strings"

	living_terminal "github.com/jwpkg/living-terminal"
	"golang.org/x/term"
)

type ProgressBar struct {
	Min             int
	Max             int
	Current         int
	Size            int
	reRenderChannel chan bool
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Min:     0,
		Max:     100,
		Current: 0,
	}
}

func (progressBar *ProgressBar) Init(requestReRender chan bool) {
	progressBar.reRenderChannel = requestReRender
}

func (progressBar *ProgressBar) Update(current int) {
	progressBar.Current = current
	progressBar.reRenderChannel <- true
}

func (progressBar *ProgressBar) SetRange(min, max int) {
	progressBar.Min = min
	progressBar.Max = max
	progressBar.reRenderChannel <- true
}

func (progressBar *ProgressBar) SetSize(size int) {
	progressBar.Size = size
	progressBar.reRenderChannel <- true
}

func (progressBar *ProgressBar) Finish() {
}

func (progressBar *ProgressBar) Render() string {
	size := progressBar.Size
	if size == 0 {
		size, _, _ = term.GetSize(int(living_terminal.OriginalStdout.Fd()))
	}
	progressRange := float64(progressBar.Max - progressBar.Min)
	progress := float64(progressBar.Current) / progressRange
	numberFilled := int(math.Max(0, math.Min(float64(size-2), math.Round(progress*float64(size-2)))))
	numberNotFilled := int(math.Max(0, math.Min(float64(size-2), float64(size-numberFilled-2))))

	return fmt.Sprintf("[%s%s]", strings.Repeat("=", numberFilled), strings.Repeat(" ", numberNotFilled))
}

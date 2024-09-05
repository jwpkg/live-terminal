package internal

import (
	"os"
)

type LineCounter struct {
	stdinCounter  *InputLineCounter
	stdoutCounter *OutputLineCounter
	stderrCounter *OutputLineCounter
}

func NewLineCounter() *LineCounter {
	lineCounter := LineCounter{
		stdinCounter:  NewInputLineCounter(os.Stdin),
		stdoutCounter: NewOuputLineCounter(os.Stdout),
		stderrCounter: NewOuputLineCounter(os.Stderr),
	}

	os.Stdin = lineCounter.stdinCounter.Reader
	os.Stdout = lineCounter.stdoutCounter.Writer
	os.Stderr = lineCounter.stderrCounter.Writer

	return &lineCounter
}

func (lineCounter *LineCounter) Count() int {
	lineCounter.pauze()
	result := lineCounter.stdoutCounter.Count() + lineCounter.stderrCounter.Count() + lineCounter.stdinCounter.Count()

	lineCounter.resume()
	return result
}

func (lineCounter *LineCounter) Stop() {
	os.Stdout = lineCounter.stdoutCounter.Stream
	os.Stderr = lineCounter.stderrCounter.Stream

	lineCounter.stdoutCounter.Stop()
	lineCounter.stderrCounter.Stop()
}

func (lineCounter *LineCounter) pauze() {
	lineCounter.stdoutCounter.preparePauze()
	lineCounter.stderrCounter.preparePauze()

	lineCounter.stdinCounter.pauze()

	os.Stdin = lineCounter.stdoutCounter.newReader
	os.Stdout = lineCounter.stdoutCounter.newWriter
	os.Stderr = lineCounter.stderrCounter.newWriter

	lineCounter.stdoutCounter.pauze()
	lineCounter.stderrCounter.pauze()
}

func (lineCounter *LineCounter) resume() {
	lineCounter.stdinCounter.resume()
	lineCounter.stdoutCounter.resume()
	lineCounter.stderrCounter.resume()
}

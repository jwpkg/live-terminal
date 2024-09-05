package internal

import (
	"bufio"
	"os"
	"strings"

	"github.com/muesli/cancelreader"
)

// Counts line on input stream
// Should be created with NewInputLineCounter()
type InputLineCounter struct {
	Stream *os.File

	Reader *os.File
	writer *os.File

	runnerDone   chan interface{}
	currentLines chan int

	cancelReader cancelreader.CancelReader
	termWatcher  TermWatcher

	count int
}

func (lineCounter *InputLineCounter) run() {
	lineCounter.cancelReader, _ = cancelreader.NewReader(lineCounter.Stream)

	scanner := bufio.NewScanner(lineCounter.cancelReader)
	scanner.Split(bufio.ScanBytes)

	go func() {
		for scanner.Scan() {

			b := scanner.Bytes()

			if len(b) >= 0 && lineCounter.termWatcher.StdinEcho {
				lineCounter.count += strings.Count(string(b), "\n")
				lineCounter.writer.Write(b)
			}
		}

		lineCounter.runnerDone <- true
	}()
}

func (lineCounter *InputLineCounter) Count() int {
	return lineCounter.count
}

func NewInputLineCounter(stream *os.File) *InputLineCounter {
	lineCounter := InputLineCounter{
		Stream:      stream,
		termWatcher: *NewTermWatcher(),
	}
	var err error

	lineCounter.Reader, lineCounter.writer, err = os.Pipe()

	if err != nil {
		panic("Cannot create capturing pipes for stdout")
	}

	lineCounter.runnerDone = make(chan interface{}, 1)
	lineCounter.currentLines = make(chan int, 1)

	lineCounter.termWatcher.Start()
	lineCounter.run()

	return &lineCounter
}

func (lineCounter *InputLineCounter) Stop() {
	lineCounter.termWatcher.Stop()
	lineCounter.pauze()
	lineCounter.writer.Close()
	lineCounter.Reader.Close()
}

func (lineCounter *InputLineCounter) pauze() {
	if lineCounter.cancelReader != nil {
		lineCounter.cancelReader.Cancel()
	}
	<-lineCounter.runnerDone
}

func (lineCounter *InputLineCounter) resume() {
	lineCounter.run()
}

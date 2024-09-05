package internal

import (
	"bufio"
	"os"
	"strings"
)

type OutputLineCounter struct {
	Stream *os.File

	reader *os.File
	Writer *os.File

	newReader *os.File
	newWriter *os.File

	runnerDone chan bool

	count int
}

func (lineCounter *OutputLineCounter) run(started chan bool) {
	scanner := bufio.NewScanner(lineCounter.reader)
	scanner.Split(bufio.ScanBytes)
	started <- true
	for scanner.Scan() {

		b := scanner.Bytes()

		if len(b) >= 0 {
			lineCounter.count += strings.Count(string(b), "\n")
			lineCounter.Stream.Write(b)
		}
	}
	lineCounter.runnerDone <- true
}

func (lineCounter *OutputLineCounter) Count() int {
	return lineCounter.count
}

func NewOuputLineCounter(stream *os.File) *OutputLineCounter {
	lineCounter := OutputLineCounter{
		Stream: stream,
	}
	var err error

	lineCounter.reader, lineCounter.Writer, err = os.Pipe()

	if err != nil {
		panic("Cannot create capturing pipes for stdout")
	}

	lineCounter.runnerDone = make(chan bool)

	started := make(chan bool, 1)
	go lineCounter.run(started)
	<-started

	return &lineCounter
}

func (lineCounter *OutputLineCounter) Stop() {
	lineCounter.Writer.Close()

	<-lineCounter.runnerDone
	lineCounter.reader.Close()
}

func (lineCounter *OutputLineCounter) preparePauze() {
	var err error
	lineCounter.newReader, lineCounter.newWriter, err = os.Pipe()

	if err != nil {
		panic("Cannot create capturing pipes for stdout")
	}

}

func (lineCounter *OutputLineCounter) pauze() {
	lineCounter.Stop()

	lineCounter.reader = lineCounter.newReader
	lineCounter.Writer = lineCounter.newWriter
}

func (lineCounter *OutputLineCounter) resume() {
	started := make(chan bool, 1)
	go lineCounter.run(started)
	<-started
}

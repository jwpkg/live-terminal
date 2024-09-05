package main

import (
	"time"
)

func main() {
	sinusWave := NewSinusWave()
	defer sinusWave.Finish()

	time.Sleep(time.Second * 10)
}

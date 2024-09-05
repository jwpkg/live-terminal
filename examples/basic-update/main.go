package main

import (
	"fmt"
	"time"

	living_terminal "github.com/living-terminal"
)

func main() {
	fmt.Println("Hello")

	line := living_terminal.NewLine("Live")
	defer line.Finish()

	done := make(chan bool)
	go func() {
		for i := range 10 {
			line.Update(fmt.Sprint("Live ", i))
			time.Sleep(time.Millisecond * 500)
		}
		done <- true
	}()
	fmt.Println("Just print below")
	<-done
}

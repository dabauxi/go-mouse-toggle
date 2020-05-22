package main

import (
	"flag"
	"fmt"
	"github.com/go-vgo/robotgo"
	_"github.com/go-vgo/robotgo"
	"time"
)

func main() {
	stepSize := 10
	fmt.Println("Starting mouse toggle.")
	duration := flag.Int("duration", 0, "Provide duration in minutes. (default 0 = infinite)")
	flag.Parse()
	_, maxY := robotgo.GetScreenSize()
	if *duration == 0 {
		fmt.Println("Duration: infinite")
		for {
			toggle(maxY, stepSize)
		}
	} else {
		fmt.Printf("Duration: %d minutes\n", *duration)
		seconds := *duration * 60
		for i := 0; i < seconds; i = i+2 {
			toggle(maxY, stepSize)
		}
	}
	fmt.Println("Mouse toggle session ended.")
}

func toggle(maxY, stepSize int) {
	toggleDown(maxY, stepSize)
	time.Sleep(1 * time.Second)
	toggleUp(maxY, stepSize)
	time.Sleep(1 * time.Second)
}

func toggleUp(maxY, stepSize int) {
	x, y := robotgo.GetMousePos()
	if (y - stepSize) <= stepSize {
		toggleDown(maxY, stepSize)
	} else {
		robotgo.DragMouse(x, y-stepSize)
	}
}

func toggleDown(maxY, stepSize int) {
	x, y := robotgo.GetMousePos()
	if (y + stepSize) >= maxY {
		toggleUp(maxY, stepSize)
	} else {
		robotgo.DragMouse(x, y+stepSize)
	}
}


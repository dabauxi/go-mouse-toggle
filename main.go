package main

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/go-vgo/robotgo"
	_ "github.com/go-vgo/robotgo"
	"strconv"
	"time"
)

/*
 compiling troubles: https://stackoverflow.com/questions/21322707/zlib-header-not-found-when-cross-compiling-with-mingw
 */

func main() {
	app := app.New()
	duration := 0
	w := app.NewWindow("Go mouse toggle")
	isRunning := false
	input := widget.NewEntry()
	input.SetPlaceHolder("infinite")
	errorLabel := widget.NewLabel("Wrong input, provide minutes.")
	errorLabel.Hide()
	progressBar := widget.NewLabel("Running 0/x minutes.")
	infiniteProgressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()
	infiniteProgressBar.Hide()
	var quitChannel chan bool = nil
	w.SetContent(widget.NewVBox(
		progressBar,
		infiniteProgressBar,
		widget.NewLabel("Enter minutes:"),
		errorLabel,
		input,
		widget.NewButton("Start", func() {
			quitChannel = make(chan bool)
			inputDuration := input.Text
			i := 0
			var err error = nil
			if inputDuration != "" {
				i, err = strconv.Atoi(input.Text)
			}
			if err != nil {
				errorLabel.Show()
				input.SetText("")
			} else {
				if isRunning {
					return
				}
				errorLabel.Hide()
				duration = i
				isRunning = true
				go func() {
					if duration == 0 {
						infiniteProgressBar.Show()
					} else {
						infiniteProgressBar.Hide()
						progressBar.Show()
						for i := 0; i < duration; i++ {
							select {
								case <- quitChannel:
									break
								default:
									labelText := fmt.Sprintf("%d minutes left running.", i + 1)
									progressBar.SetText(labelText)
									time.Sleep(time.Second * 60)
								}
						}
					}
				}()
				go runApp(duration, quitChannel)
			}
		}),
		widget.NewButton("Stop", func() {
			if isRunning {
				quitChannel <- true
				infiniteProgressBar.Hide()
				progressBar.Hide()
				close(quitChannel)
				isRunning = false
			}
		}),
		widget.NewButton("Quit", func() {
			if isRunning {
				quitChannel <- true
				close(quitChannel)
			}
			app.Quit()
		}),
	))
	w.ShowAndRun()
}

func runApp(duration int, quit chan bool) {
	stepSize := 10
	_, maxY := robotgo.GetScreenSize()
	if duration == 0 {
		for {
			select {
			case <- quit:
				return
			default:
				toggle(maxY, stepSize)
			}
		}
	} else {
		seconds := duration * 60
		for i := 0; i < seconds; i = i+2 {
			select {
			case <-quit:
				return
			default:
				toggle(maxY, stepSize)
			}
		}
	}
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


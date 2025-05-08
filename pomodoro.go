package pomodoro

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func Run() {
	userStop := setupUserStop()
	pomodoro := newTimer("🍅 Pomodoro", 25*time.Minute, userStop)
	shortBreak := newTimer("🌱 Short break", 5*time.Minute, userStop)
	longBreak := newTimer("🌴 Long break", 15*time.Minute, userStop)
	var pomodoroCount int

	getNextBreak := func() timer {
		pomodoroCount++
		if pomodoroCount < 4 {
			return shortBreak
		}

		pomodoroCount = 0
		return longBreak
	}

	for {
		if err := pomodoro(); err != nil {
			break
		}

		nextBreak := getNextBreak()
		if err := nextBreak(); err != nil {
			break
		}
	}
	clearScreen()
	fmt.Print("Goodbye! You've ended the Pomodoro session. 🌿\n")
}

func setupUserStop() (userStop chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	userStop = make(chan struct{})
	go func() {
		<-signals
		close(userStop)
	}()
	return userStop
}

type timer func() error

func newTimer(
	label string,
	duration time.Duration,
	userStop chan struct{},
) timer {
	return func() error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		end := time.Now().Add(duration)

		clearScreen()
		fmt.Printf("%s %02d:%02d",
			label,
			int(duration.Minutes()),
			int(duration.Seconds())%60,
		)

		for {
			select {
			case <-ticker.C:
				remaining := end.Sub(time.Now()).Round(time.Second)
				if remaining <= 0 {
					return nil
				}

				clearScreen()
				fmt.Printf("%s %02d:%02d",
					label,
					int(remaining.Minutes()),
					int(remaining.Seconds())%60,
				)

			case <-userStop:
				return errors.New("stopped by user")
			}
		}
	}
}

func clearScreen() {
	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H")  // move cursor to top-left
}

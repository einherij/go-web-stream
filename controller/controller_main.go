package controller

import (
	"fmt"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// pins on GPIO
	pinl1 = rpio.Pin(5)
	pinl2 = rpio.Pin(6)
	pinr1 = rpio.Pin(13)
	pinr2 = rpio.Pin(19)
)

// StartRPIO prepares pins for work
func StartRPIO() {
	// Start rpio
	err := rpio.Open()
	if err != nil {
		fmt.Println("Error starting RPIO, error:", err)
	}
	pinl1.Output()
	pinl2.Output()
	pinr1.Output()
	pinr2.Output()
	defer pinl1.Low()
	defer pinl2.Low()
	defer pinr1.Low()
	defer pinr2.Low()
}

// StopRPIO disables RPIO all pins
func StopRPIO() {
	rpio.Close()
}

// RunRobot sends commands to Robot
func RunRobot(r int) {
	switch r {
	case 1:
		startRForward()
	case 11:
		stopRightSide()
	case 2:
		startGoForward()
	case 22:
		stopMoving()
	case 3:
		startLForward()
	case 33:
		stopLeftSide()
	case 4:
		startGoRight()
	case 44:
		stopMoving()
	case 6:
		startGoLeft()
	case 66:
		stopMoving()
	case 7:
		startLBackward()
	case 77:
		stopLeftSide()
	case 8:
		startGoBackward()
	case 88:
		stopMoving()
	case 9:
		startRBackward()
	case 99:
		stopRightSide()
	default:
		stopMoving()
	}
}

func startLBackward() {
	pinl1.High()
	pinl2.Low()
}

func stopLeftSide() {
	pinl1.Low()
	pinl2.Low()
}

func startLForward() {
	pinl1.Low()
	pinl2.High()
}

func startRForward() {
	pinr1.High()
	pinr2.Low()
}

func stopRightSide() {
	pinr1.Low()
	pinr2.Low()
}

func startRBackward() {
	pinr1.Low()
	pinr2.High()
}

func startGoForward() {
	startLForward()
	startRForward()
}

func startGoBackward() {
	startLBackward()
	startRBackward()
}

func stopMoving() {
	stopLeftSide()
	stopRightSide()
}

func startGoLeft() {
	startLForward()
	startRBackward()
}

func startGoRight() {
	startLBackward()
	startRForward()
}

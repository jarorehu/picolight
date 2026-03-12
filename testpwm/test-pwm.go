package main

import (
	"machine"
	"slices"
	"time"
)

// ----------------------------------------
// https://github.com/tinygo-org/tinygo/issues/2583
type PWM interface {
	Set(channel uint8, value uint32)
	SetPeriod(period uint64) error
	Enable(bool)
	Top() uint32
	Configure(config machine.PWMConfig) error
	Channel(machine.Pin) (uint8, error)
}

type config struct {
	pinLed    machine.Pin
	percent   uint32
	generator PWM
	channel   uint8
}

var period uint64 = 1e9 / 2000

var ledconfig = map[string]config{
	"R": config{percent: 50, generator: machine.PWM6, pinLed: machine.GPIO12, channel: 0},
	"G": config{percent: 50, generator: machine.PWM4, pinLed: machine.GPIO8, channel: 0},
	"B": config{percent: 50, generator: machine.PWM2, pinLed: machine.GPIO4, channel: 0},
	"W": config{percent: 50, generator: machine.PWM0, pinLed: machine.GPIO0, channel: 0},
}

func blik(delay int, cnt int) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := 0; i < cnt; i++ {
		led.Low()
		time.Sleep(time.Millisecond * time.Duration(delay))

		led.High()
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func main() {

	var btnColors = [4]string{"R", "G", "B", "W"}
	//	var intensity = []int{0, 1, 2, 4, 6, 10, 16, 26, 42, 68, 100}
	var intensity = []int{100, 95, 90, 80, 70, 55, 25, 10, 5, 0} // inverzni stupnice

	// PWM diody, nastaveni pinu a pwm generátoru
	//var pwm PWM
	var ch uint8
	for col, led := range ledconfig {
		led.pinLed.Configure(machine.PinConfig{Mode: machine.PinPWM})
		//pwm, ch, _ = getPWM(led.pinLed)
		//led.generator = pwm
		led.generator.Configure(machine.PWMConfig{Period: period})
		ch, _ = led.generator.Channel(led.pinLed)
		led.generator.Set(ch, led.generator.Top())
		led.channel = ch
		ledconfig[col] = led
	}
	//machine.PWMPeripheral(0)
	println("configured!", ledconfig)
	blik(100, 3)

	for {
		// color lights UP
		for _, col := range btnColors {
			led := ledconfig[col]
			for _, i := range intensity {
				led.generator.Set(led.channel, led.generator.Top()*uint32(i)/100)
				blik(200, 1)
				time.Sleep(200 * time.Millisecond)
			}
		}
		// color lights DOWN
		for _, col := range btnColors {
			led := ledconfig[col]
			for _, i := range slices.Backward(intensity) {
				led.generator.Set(led.channel, led.generator.Top()*uint32(i)/100)
				blik(200, 1)
				time.Sleep(200 * time.Millisecond)
			}
		}

	}

}

package main

import (
	"machine"
	"time"
)

const (
	timeDebounce     = 10 * time.Millisecond
	timeLongPress    = 500 * time.Millisecond
	timeDblClickWait = 250 * time.Millisecond
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

/*
func getPWM(pin machine.Pin) (PWM, uint8, error) {
	var pwms = [...]PWM{machine.PWM0, machine.PWM1, machine.PWM2, machine.PWM3, machine.PWM4, machine.PWM5, machine.PWM6, machine.PWM7}
	slice, err := machine.PWMPeripheral(pin)
	if err != nil {
		return nil, 0, err
	}
	pwm := pwms[slice]
	err = pwm.Configure(machine.PWMConfig{Period: 1e9 / 100}) // 100Hz for starters.
	if err != nil {
		return nil, 0, err
	}
	channel, err := pwm.Channel(pin)
	return pwm, channel, err
}
*/
// ----------------------------------------

type config struct {
	pinLed    machine.Pin
	pinPower  int
	generator PWM
	channel   uint8
}

var period uint64 = 1e9 / 2000

var btnconfig = map[string]machine.Pin{
	"R":     machine.GPIO11,
	"G":     machine.GPIO10,
	"B":     machine.GPIO3,
	"W":     machine.GPIO2,
	"plus":  machine.GPIO15,
	"minus": machine.GPIO14,
}

var ledconfig = map[string]config{
	"R": config{pinPower: 1, generator: machine.PWM6, pinLed: machine.GPIO12, channel: 0},
	"G": config{pinPower: 1, generator: machine.PWM4, pinLed: machine.GPIO8, channel: 0},
	"B": config{pinPower: 1, generator: machine.PWM2, pinLed: machine.GPIO4, channel: 0},
	"W": config{pinPower: 1, generator: machine.PWM0, pinLed: machine.GPIO0, channel: 0},
}

var pwmScale = []uint32{100, 95, 90, 80, 70, 55, 25, 10, 5, 0} // inverzni stupnice

func buttonAction(btn machine.Pin) string {
	action := ""
	pressTime := time.Now()

	time.Sleep(timeDebounce) // sure for change
	if !btn.Get() {
		action = "short"
		for !btn.Get() {
			time.Sleep(timeDebounce)
		}
		if time.Since(pressTime) > timeLongPress {
			// long time achieved
			action = "long"
			time.Sleep(timeDebounce)
		} else {
			time.Sleep(timeDebounce) // sure for change
			// wait for Dbl click pause
			releaseTime := time.Now()
			for btn.Get() && time.Since(releaseTime) < timeDblClickWait {
				time.Sleep(timeDebounce)
			}
			if !btn.Get() {
				action = "double"
				time.Sleep(timeDebounce)
				for !btn.Get() {
					time.Sleep(timeDebounce)
				}
			}
		}
	}

	return action
}

func main() {

	var btnColors = [4]string{"R", "G", "B", "W"}
	var btnAction = [2]string{"plus", "minus"}
	var released bool = false
	var selectedColor string = "all"
	var selectedAction string = ""

	// tlacitka, definice a nastaveni
	for _, pin := range btnconfig {
		pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// PWM diody, nastaveni pinu a pwm generátoru
	//var pwm PWM
	var ch uint8
	for col, led := range ledconfig {
		led.pinLed.Configure(machine.PinConfig{Mode: machine.PinPWM})
		led.generator.Configure(machine.PWMConfig{Period: period})
		ch, _ = led.generator.Channel(led.pinLed)
		led.generator.Set(ch, led.generator.Top()*pwmScale[led.pinPower]/100)
		led.channel = ch
		ledconfig[col] = led
	}
	//machine.PWMPeripheral(0)
	println("configured!", ledconfig)

	for {
		// color selection buttons
		for _, col := range btnColors {
			if !btnconfig[col].Get() {
				// pressed !!
				time.Sleep(25 * time.Millisecond)
				if !btnconfig[col].Get() {
					// after 25 ms still pressed
					if selectedColor == col {
						selectedColor = "all"
					} else {
						selectedColor = col
					}
				}
				// wait for release
				released = false
				for !released {
					released = true
					time.Sleep(25 * time.Millisecond)
					if !btnconfig[col].Get() {
						released = false
					}
				}
			}
		}

		// change up down buttons
		for _, act := range btnAction {
			if !btnconfig[act].Get() {
				// pressed !!
				time.Sleep(25 * time.Millisecond)
				if !btnconfig[act].Get() {
					// after 25 ms still pressed
					selectedAction = act
				} else {
					selectedAction = ""
				}
				// wait for release
				released = false
				for !released {
					released = true
					time.Sleep(25 * time.Millisecond)
					if !btnconfig[act].Get() {
						released = false
					}
				}
			}
		}
		// debug result
		println(selectedColor, selectedAction)

		if selectedAction != "" {
			for col, led := range ledconfig {
				if selectedColor == "all" || selectedColor == col {
					// increment/decrement
					switch selectedAction {
					case "plus":
						led.pinPower += 1
						if led.pinPower > len(pwmScale) {
							led.pinPower = len(pwmScale)
						}
					case "minus":
						led.pinPower -= 1
						if led.pinPower < 0 {
							led.pinPower = 0
						}
					}
					// set pwm
					led.generator.Set(led.channel, led.generator.Top()*pwmScale[led.pinPower]/100)

				}
			}
			selectedAction = ""
		}

		time.Sleep(200 * time.Millisecond)

	}

}

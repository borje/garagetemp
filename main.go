package main

import (
	"fmt"
	"time"

	"github.com/felixge/pidctrl"
)

type system struct {
	dropoffPerMinute float64
	temp             float64
}

func (s *system) Update(warmup float64) {
	s.temp = s.temp + warmup - s.dropoffPerMinute
}

func (s *system) GetTemp() float64 {
	return s.temp
}

func main() {
	garage := system{.025, 19}
	//	c := pidctrl.NewPIDController(20, .1, 0)
	c := pidctrl.NewPIDController(20, .1, 0)
	c.Set(20)
	c.SetOutputLimits(0, 5)
	duration := time.Minute * 5
	for {
		output := c.UpdateDuration(garage.GetTemp(), duration)
		garage.Update(output * .1 / 5 * 2000 / 750)
		fmt.Printf("Garage: %.2f\tWarmup: %.1f\n", garage.GetTemp(), output)
		time.Sleep(time.Second)
	}

}

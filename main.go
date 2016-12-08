package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/felixge/pidctrl"
	chart "github.com/wcharczuk/go-chart"
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
	//
	file, err := os.Create("output.png")
	if err != nil {
		log.Fatal("Cant create file")
	}
	const size = 12 * 4
	//xvalues := make([]float64, size, size)
	xvalues := make([]time.Time, size, size)
	yvalues := make([]float64, size, size)
	/*
		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},
			},
		}
	*/

	garage := system{.025, 19}
	//	c := pidctrl.NewPIDController(20, .1, 0)
	c := pidctrl.NewPIDController(20, .05, 0)
	c.Set(20)
	c.SetOutputLimits(0, 5)
	duration := time.Minute * 5
	theTime := time.Now()
	for i := 0; i < size; i++ {
		//xvalues[i] = float64(i)
		xvalues[i] = theTime
		yvalues[i] = garage.GetTemp()

		output := c.UpdateDuration(garage.GetTemp(), duration)
		garage.Update(output * .1 / 5)
		fmt.Printf("Garage: %.2f\tWarmup: %.1f\n", garage.GetTemp(), output)
		//time.Sleep(time.Second)
		theTime = theTime.Add(duration)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			TickPosition:   chart.TickPositionBetweenTicks,
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xvalues,
				YValues: yvalues,
			},
		},
	}

	err = graph.Render(chart.PNG, file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

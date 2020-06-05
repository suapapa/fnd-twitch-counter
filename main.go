package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/suapapa/go_devices/tm1638"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var (
	targetCnt int
	loopSecs  int
)

func init() {
	flag.IntVar(&targetCnt, "t", 50, "target count")
	flag.IntVar(&loopSecs, "l", 30, "update every given senconds")
}

func main() {
	flag.Parse()

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	dev, err := tm1638.Open(
		gpioreg.ByName("17"), // data
		gpioreg.ByName("27"), // clk
		gpioreg.ByName("22"), // stb
	)
	if err != nil {
		panic(err)
	}
	displayWelcome(dev)

	updateFND := func(dev *tm1638.Module) int {
		cnt := getFollowerCnt()
		cntStr := []rune(fmt.Sprintf("%4d%4d", cnt, targetCnt))
		for i := 0; i < 8; i++ {
			var dot bool
			if (i+1)%4 == 0 {
				dot = true
			}
			dev.SetChar(i, cntStr[i], dot)
		}

		// fill bar graph leds
		currLedCnt := cnt * 8 / targetCnt
		if cnt >= targetCnt {
			for i := 0; i < 8; i++ {
				dev.SetLed(i, tm1638.Red)
			}
		} else {
			for i := 0; i < 8; i++ {
				if i < currLedCnt {
					dev.SetLed(8-i-1, tm1638.Green)
				} else {
					dev.SetLed(8-i-1, tm1638.Off)
				}
			}
		}
		return cnt
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	ticker := time.NewTicker(time.Duration(loopSecs) * time.Second)
loop:
	for {
		select {
		case t := <-ticker.C:
			cnt := updateFND(dev)
			log.Printf("%v: %d", t, cnt)
		case <-quit:
			break loop
		}
	}
	log.Println("byebye")
}

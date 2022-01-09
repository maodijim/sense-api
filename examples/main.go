package main

import (
	"fmt"
	"github.com/maodijim/sense-api"
	"time"
)

func main() {
	s, _ := sense.NewSenseApi("test@test.com", "test")
	timer := time.After(time.Second * 300)

	// Get always on
	al, _ := s.AlwaysOn()

	// Get device event timeline
	tl, _ := s.TimeLine(30)

	// Get list of devices
	do, _ := s.DevicesOverview(true)

	// Get Trend
	start := time.Now().Add(time.Hour * time.Duration(24*-30))
	t, _ := s.Trend(sense.TrendMonth, start)
	t, _ = s.Trend(sense.TrendWeek, start)
	t, _ = s.Trend(sense.TrendDay, start)
	t, _ = s.Trend(sense.TrendYear, start)
	// Renew access token
	_ = s.RenewToken()

	fmt.Printf("%v\n%v\n%v\n%v", al, tl, do, t)
	for {
		select {
		case <-timer:
			return
		default:
			msg, err := s.ReadMessage()
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			fmt.Println(msg)
			time.Sleep(3 * time.Second)
		}
	}
}

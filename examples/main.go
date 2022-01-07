package main

import (
	"fmt"
	sense_api "github.com/maodijim/sense"
	"time"
)

func main() {
	s, _ := sense_api.NewSenseApi("test@test.com", "test")
	timer := time.After(time.Second * 300)

	// Get always on
	al, _ := s.AlwaysOn()

	// Get device event timeline
	tl, _ := s.TimeLine(30)

	// Get list of devices
	do, _ := s.DevicesOverview(true)

	// Get Trend
	start := time.Now().Add(time.Hour * time.Duration(24*-30))
	t, _ := s.Trend(sense_api.TrendMonth, start)
	t, _ = s.Trend(sense_api.TrendWeek, start)
	t, _ = s.Trend(sense_api.TrendDay, start)
	t, _ = s.Trend(sense_api.TrendYear, start)
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

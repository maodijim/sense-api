package main

import (
	"fmt"
	"time"

	"github.com/maodijim/sense-api"
)

func main() {
	s, _ := sense.NewSenseApi("test@test.com", "test")

	// Get Time of Use Rate Zones
	rz, _ := s.RateZone()

	// Get history comparisons
	hc, _ := s.GetHistoryComparison()

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

	fmt.Printf("%v\n%v\n%v\n%v\n%v\n%v\n", al, tl, do, t, hc, rz)

	shouldClose := make(chan bool, 1)
	go func() {
		_ = s.ReadMessageAsync(shouldClose)
	}()
	time.Sleep(time.Second * 5)
	msgs, _ := s.ReadMessages()
	shouldClose <- true
	fmt.Printf("%d messages\n", len(msgs))
	err := s.Close()
	fmt.Println(err)
	msg, err := s.ReadMessage()
	fmt.Println(msg)

	timer := time.After(time.Second * 5)
	for {
		select {
		case <-timer:
			return
		default:
			msg, err := s.ReadMessage()
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			fmt.Printf("message: %s\n", msg)
		}
	}
}

package api

import (
	"log"
	"time"
	"urlshortener/pkg/protobuf/apipb"
)

func schedulePeriodicTask(task func(), d time.Duration) (stop func()) {
	ticker := time.NewTicker(d)
	done := make(chan struct{}, 1)
	go func(ticker *time.Ticker, task func(), done <-chan struct{}) {
		for {
			select {
			case <-ticker.C:
				task()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}(ticker, task, done)
	return func() { done <- struct{}{} }
}

func makePingCall(client pingCallable, name string, active *bool) {
	c, cancel := makeCtx()
	defer cancel()
	_, err := client.Ping(c, &apipb.PingRequest{})
	if err != nil {
		*active = false
		log.Printf("failed to ping %s: %v", name, err)
	} else {
		*active = true
		log.Printf("%s is active", name)
	}
}

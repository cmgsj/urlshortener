package api

import (
	"log"
	"time"
	"urlshortener/pkg/protobuf/apipb"
	"urlshortener/pkg/protobuf/cachepb"
	"urlshortener/pkg/protobuf/urlspb"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	urlClient          urlspb.UrlsClient
	urlServiceActive   bool
	cacheClient        cachepb.CacheClient
	cacheServiceActive bool
	router             *gin.Engine
	trustedProxies     []string
}

func (server *httpServer) registerTrustedProxies() {
	err := server.router.SetTrustedProxies(server.trustedProxies)
	if err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
}

func (server *httpServer) pingServices() {
	c, cancel := makeCtx()
	defer cancel()
	_, err := server.urlClient.Ping(c, &apipb.PingRequest{})
	if err != nil {
		server.urlServiceActive = false
		log.Println("failed to ping url service: ", err)
	} else {
		server.urlServiceActive = true
		log.Println("url service is active")
	}
	c, cancel = makeCtx()
	defer cancel()
	_, err = server.cacheClient.Ping(c, &apipb.PingRequest{})
	if err != nil {
		server.cacheServiceActive = false
		log.Println("failed to ping cache service: ", err)
	} else {
		server.cacheServiceActive = true
		log.Println("cache service is active")
	}
}

func (server *httpServer) schedulePeriodicTask(task func(), d time.Duration) func() {
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

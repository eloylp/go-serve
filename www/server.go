package www

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Shutdown(s *http.Server, timeout time.Duration) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
}

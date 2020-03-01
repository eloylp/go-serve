package www

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown(ctx context.Context, s *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		if err := s.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
}

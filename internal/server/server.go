package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const listenAddr = "localhost:8080"

func Start() {
	srv := http.Server{Addr: listenAddr, Handler: timeoutMiddleware(serverHandler())}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("can't start server: %s", err)
		}
	}()

	log.Printf("server started at: %s", listenAddr)

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigs:
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("can't shutdown server: %s", err)
		}
		log.Println("server stopped")
	}
}

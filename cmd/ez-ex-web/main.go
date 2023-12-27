package main

import (
	"github.com/armanimichael/ez-ex/cmd/ez-ex-web/controller"
	"github.com/armanimichael/ez-ex/cmd/ez-ex-web/httpserver"
	"github.com/armanimichael/ez-ex/cmd/ez-ex-web/service"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	// TODO: use config files
	var err error
	var authService controller.Authenticator = service.NewAuthService(64, 1, 65536, uint8(runtime.NumCPU()))

	handler := controller.New(authService)
	httpServer := httpserver.Start(handler, httpserver.WithPort("1024"))
	defer func() {
		if err = httpServer.Shutdown(); err != nil {
			log.Fatalf("server shutdown error: %v", err)
		}
		log.Println("server shut down successfully")
	}()

	// Waiting exit signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	select {
	case sig := <-interrupt:
		log.Println("server interrupt signal: " + sig.String())
	case err = <-httpServer.Notify():
		log.Fatalf("server error: %v", err)
	}
}

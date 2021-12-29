package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EstebanFallaGlobant/gorillaMux-demo/api"
)

func main() {
	var dir string

	flag.StringVar(&dir, "addr", "127.0.0.1:8080", "Addres in which the server runs")
	flag.Parse()

	logger := log.Default()
	api := api.API{Logger: logger}
	svr := &http.Server{
		Addr:         dir,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.Initialize(),
	}

	go func() {
		logger.Printf("Server starting at: %s\n", svr.Addr)
		if msg, err := api.GetHealthURL(); err == nil {
			logger.Printf("For health information: %s%s\n", dir, msg)
		}

		if err := svr.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	svr.Shutdown(ctx)

	logger.Println("Shutting down")

	os.Exit(0)
}

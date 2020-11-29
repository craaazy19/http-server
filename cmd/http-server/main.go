package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/craaazy19/http-server/internal/routing"
	"golang.org/x/net/netutil"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	router := &routing.Router{}

	http.HandleFunc("/parse", router.Parse)
	http.HandleFunc("/test", router.Test)

	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = l.Close() }()

	l = netutil.LimitListener(l, 100)

	srv := &http.Server{
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err = srv.Serve(l)
		if err != nil {
			log.Print(err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	router.Shutdown()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%s", err)
	}

	log.Printf("server exited properly")
}

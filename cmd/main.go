package main

import (
	"context"
	"github.com/ivanovaleksey/simdrone/dispatcher"
	"github.com/ivanovaleksey/simdrone/storage"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't get wd"))
	}

	go http.ListenAndServe(":8080", nil)

	ctx := context.Background()

	store := storage.New(filepath.Join(wd, "storage", "data"))
	disp := dispatcher.New(store, store)

	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGTERM)

	log.Printf("PID %d\n", os.Getpid())
	log.Println("starting")
	if err := disp.Start(ctx); err != nil {
		log.Fatal(errors.Wrap(err, "dispatcher error"))
	}
	log.Println("started")

	<-sign
	log.Println("closing")

	done := make(chan error)
	go func() {
		defer close(done)
		done <- disp.Close()
	}()

	const closeTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, closeTimeout)
	defer cancel()

	var closeErr error
	select {
	case closeErr = <-done:
	case <-ctx.Done():
		closeErr = ctx.Err()
	}
	if closeErr == nil {
		log.Println("closed")
	} else {
		log.Printf("closed with err %v\n", closeErr)
	}
}

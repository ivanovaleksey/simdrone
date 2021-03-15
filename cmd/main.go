package main

import (
	"context"
	"github.com/ivanovaleksey/simdrone/dispatcher"
	"github.com/ivanovaleksey/simdrone/storage"
	"github.com/pkg/errors"
	"io"
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
	done := make(chan error)
	go func() {
		defer close(done)
		log.Println("starting")
		if err := disp.Start(ctx); err != nil {
			log.Printf("dispatcher error %v\n", err)
			done <- err
		}
	}()

	select {
	case <-sign:
	case <-done:
	}

	if err := closeApp(ctx, disp); err != nil {
		log.Printf("closed with err %v\n", err)
	}
	log.Println("closed")
}

func closeApp(ctx context.Context, app io.Closer) error {
	const closeTimeout = 5 * time.Second
	done := make(chan error)

	go func() {
		defer close(done)
		log.Println("closing")
		done <- app.Close()
	}()

	ctx, cancel := context.WithTimeout(ctx, closeTimeout)
	defer cancel()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

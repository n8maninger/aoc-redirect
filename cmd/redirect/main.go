package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	l, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	loc := time.FixedZone("UTC-5", -5*60*60)

	s := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.Trim(r.URL.Path, "/")
			sec, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "invalid timestamp %q", p)
				return
			}

			y, m, d := time.Unix(sec, 0).In(loc).Date()
			switch {
			case m != time.December:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "advent of code starts in December, not %s", m)
				return
			case d < 1:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "advent of code starts on December 1st, not %d", d)
				return
			case d > 25:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "advent of code ends on December 25th, not %d", d)
				return
			}

			u := fmt.Sprintf("https://adventofcode.com/%d/day/%d", y, d)
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
			log.Printf("Redirected to %d %d", y, d)
		}),
	}
	defer s.Close()

	go func() {
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
}

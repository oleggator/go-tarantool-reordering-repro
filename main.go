//go:build go_tarantool_ssl_disable

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	threads := flag.Int("t", 1, "GOMAXPROCS")
	concurrency := flag.Int("c", 0, "tarantool connector concurrency")
	flag.Parse()

	fmt.Printf("threads=%d concurrency=%d\n", *threads, *concurrency)

	runtime.GOMAXPROCS(*threads)

	dialer := tarantool.NetDialer{Address: "127.0.0.1:3301"}
	opts := tarantool.Opts{
		Concurrency: uint32(*concurrency),
	}
	conn, err := tarantool.Connect(context.Background(), dialer, opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	futures := make(chan *tarantool.Future, 1024)
	var eg errgroup.Group
	eg.Go(func() error {
		return writer(conn, futures)
	})
	eg.Go(func() error {
		return reader(futures)
	})

	if err := eg.Wait(); err != nil {
		log.Println(err)
		return
	}
}

func writer(conn *tarantool.Connection, futures chan<- *tarantool.Future) error {
	if _, err := conn.Do(tarantool.NewCallRequest("reset")).Get(); err != nil {
		return err
	}

	for seq := 0; ; seq++ {
		futures <- conn.Do(tarantool.NewCallRequest("test_func").Args([]any{seq}))
	}
}

func reader(futures <-chan *tarantool.Future) error {
	var (
		ok     = 0
		errors = 0
	)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			sum := ok + errors
			fmt.Printf("\rok: %d (%.2f%%) | out of order: %d (%.2f%%)",
				ok, float32(ok)/float32(sum)*100,
				errors, float32(errors)/float32(sum)*100,
			)

		case future := <-futures:
			resp, err := future.Get()
			if err != nil {
				return err
			}

			switch data := resp.Data[0].(string); data {
			case "ok":
				ok++
			case "err":
				errors++
			default:
				panic("invalid request" + data)
			}
		}
	}
}

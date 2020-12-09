package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	errorGroup, _ := errgroup.WithContext(ctx)

	// boot a normal server
	errorGroup.Go(func() error {
		return bootServer(ctx, 10086, "server")
	})
	// boot a debug server
	errorGroup.Go(func() error {
		return bootServer(ctx, 10010, "debug")
	})

	// listen signal
	errorGroup.Go(func() error {
		return listenSignal(ctx, cancel)
	})

	// wait err
	if err := errorGroup.Wait(); err != nil {
		fmt.Printf("%+v", err)
	}

	// wait for the server to process the request which not finished
	quiteCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	<-quiteCtx.Done()

	fmt.Println("all server stoped, program quit")
}

func bootServer(ctx context.Context, port uint, name string) error {
	fmt.Printf("[%s:%d] booting ...\n", name, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/index", httpHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		stopSever(server, name)
	}()

	err := server.ListenAndServe()

	return err
}

func stopSever(s *http.Server, name string) {
	fmt.Printf("[%s%s] : stoping ... \n", name, s.Addr)

	// graceful shutdown
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_ = s.Shutdown(ctx)
}

func listenSignal(ctx context.Context, cancel context.CancelFunc) error {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-signalChannel:
		fmt.Printf("\ngot signal : %s \n", s)

		err := errors.New("server was stoped by signal \n")
		cancel()

		return err
	case <-ctx.Done():
		return nil
	}
}

func httpHandler(rsp http.ResponseWriter, rqs *http.Request) {
	fmt.Println("http rqs", rqs.URL)

	// simulation time-consuming operation
	quiteCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	<-quiteCtx.Done()

	_, _ = fmt.Fprintf(rsp, fmt.Sprintf("request %s", rqs.URL))
}

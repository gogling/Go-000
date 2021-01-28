package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Message struct {
	listenerId int
	message    string
}

var channel = make(chan Message, 100)

func main() {
	fmt.Println("starting the server ...")

	listener, errListen := net.Listen("tcp", "127.0.0.1:3000")
	if errListen != nil {
		fmt.Println("error listening", errListen.Error())
		return
	}

	ctx, errorGroup := ctrlexist(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error:%v\n", err)
			return
		}

		logic(ctx, errorGroup, conn)
	}
}

func logic(ctx context.Context, errorGroup *errgroup.Group, conn net.Conn) {
	// boot gorutine deal conn read
	errorGroup.Go(func() error {
		return startListen(ctx, conn, 1)
	})
	errorGroup.Go(func() error {
		return startListen(ctx, conn, 2)
	})

	// boot gorutine deal chan msg
	errorGroup.Go(func() error {
		return biz(ctx, 1)
	})
	errorGroup.Go(func() error {
		return biz(ctx, 2)
	})

	fmt.Println("all server stoped, program quit")
}

func startListen(ctx context.Context, conn net.Conn, listenerId int) error {
	fmt.Printf("listener : %d : start listen ... \n", listenerId)

	for {
		buf := make([]byte, 512)
		buflen, err := conn.Read(buf)

		if err != nil {
			fmt.Println("read error or finished", err.Error())

			_ = conn.Close()

			return err
		}

		cnt := string(buf[:buflen])

		select {
		case channel <- Message{listenerId: listenerId, message: cnt}:
			fmt.Printf("listener : %d : read conn cnt : %v \n", listenerId, cnt)
		case <-ctx.Done():
			fmt.Printf("listener : %d : stop listen; \n", listenerId)

			return nil
		}
	}
}

func biz(ctx context.Context, bizId int) error {
	fmt.Printf("biz : %d : start do biz ... \n", bizId)

	for {
		select {
		case msg := <-channel:
			fmt.Printf("bizer : %d : read msg from chan : %v read by listener: %d \n", bizId, msg.message, msg.listenerId)
		case <-ctx.Done():
			fmt.Printf("biz : %d : stop do biz ... \n", bizId)

			return nil
		}
	}
}

func ctrlexist(listener net.Listener) (context.Context, *errgroup.Group) {
	ctx, cancel := context.WithCancel(context.Background())
	errorGroup, _ := errgroup.WithContext(ctx)

	// listen signal
	errorGroup.Go(func() error {
		return listenSignal(ctx, cancel)
	})

	go func() {
		<-ctx.Done()
		stopSever(listener)
	}()

	return ctx, errorGroup
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

func stopSever(listener net.Listener) {
	fmt.Printf("stoping listen... \n")

	_ = listener.Close()
}

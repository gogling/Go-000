
### 作业：
Q : 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。

##### 回答：
如下为具体实现 :
```go
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
```


执行流程如下:
```
➜  Week03 git:(main) ✗ go run main.go
# 启动2个server
[server:10086] booting ...
[debug:10010] booting ...

# 处理请求
http rqs /index
http rqs /index

# ctrl+c 发送interrupt sinal
^C
got signal : interrupt 

# 影响signal退出server
[server:10086] : stoping ... 
[debug:10010] : stoping ... 
server was stoped by signal 
all server stop, program quit
```

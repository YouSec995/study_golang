package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stop := make(chan struct{})
	g, ctx := errgroup.WithContext(context.Background())

	var exit = make(chan os.Signal)
	//监听 Ctrl+C 信号
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		ctx, _ = context.WithTimeout(ctx, 2*time.Second)

		mux := http.NewServeMux()
		mux.Handle("/", &helloGoper{})
		server := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		time.Sleep(3 * time.Second) // 假设有三秒延迟发生，超过2秒超时上限
		fmt.Println("ctx1")
		err := CheckErr(ctx, stop, exit)

		go func() {
			// 若 stop 没有被关闭时，则 goroutine 阻塞在这儿
			<-stop
			// 在 stop 关闭之后调用 shutDown 关闭这个 ListenAndServe
			server.Shutdown(context.Background())
		}()

		if err != nil {
			// do sth before quit goroutine
			return err
		}
		return server.ListenAndServe()
	})

	g.Go(func() error {
		ctx, _ = context.WithTimeout(ctx, 2*time.Second)
		mux := http.NewServeMux()
		mux.Handle("/", &helloWord{})
		server := &http.Server{
			Addr:    ":9091",
			Handler: mux,
		}

		time.Sleep(1 * time.Second)
		fmt.Println("ctx2")
		err := CheckErr(ctx, stop, exit)
		go func() {
			<-stop
			server.Shutdown(context.Background())
		}()
		if err != nil {
			// do sth
			return err
		}
		return server.ListenAndServe()
	})

	err := g.Wait()
	if err == nil {
		fmt.Println("All server closed...")
	} else {
		// main 程序可以做一些事情
		fmt.Println("get some err:", err)
	}
}

type helloGoper struct{}

type helloWord struct{}

func (*helloGoper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Goper!")
}

func (*helloWord) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// 用于接收携程的错误
func CheckErr(ctx context.Context, stop chan struct{}, exit chan os.Signal) error {
	select {
	case <-ctx.Done():
		// 若协程中有错误，则将通道 stop 关闭
		close(stop)
		return ctx.Err()
	case <-exit:
		close(stop)
		return errors.New("killed...")
	default:
		return nil
	}
}

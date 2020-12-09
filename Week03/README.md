# 学习笔记

## 作业题目

1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

## 作业内容

通过`errgroup`的源码，可以发现其`WithContext`会返回一个`WithCancel`的`Context`，并且，当`errgroup`的`Wait`方法结束后，会自动调用该`Context`的`cancel`方法。因此，当清楚了`errgroup`的基本工作原理后，本次作业即可使用`WithContext`和`Wait`结合，解决~

```go
// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
```

实现代码如下：

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
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Linux Signal
	g.Go(func() error {
		select {
		case s := <-sc:
			return errors.New(fmt.Sprintf("linux signal: receive signal [%s]", s.String()))
		case <-ctx.Done():
			fmt.Println("linux signal: shutdown by context cancel")
			return nil
		}
	})

	// HTTP Server
	g.Go(func() error {
		s := http.Server{Addr: ":10086"}

		go func() {
			<-ctx.Done()
			fmt.Println("http server: shutdown by context cancel")
			s.Shutdown(context.Background())
		}()

		return s.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("occur error in error group: %v\n", err)
	}
	fmt.Println("all server was done")
}
```


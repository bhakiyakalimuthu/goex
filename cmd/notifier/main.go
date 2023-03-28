package main

import (
	"bufio"
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bhakiyakalimuthu/goex/notifier"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	interval = flag.Duration("INTERVAL", 100*time.Millisecond, "interval in seconds")
	url      = flag.String("URL", "", "URL to post data")
)

const (
	noOfWorkers = 5
)

func main() {
	flag.Parse()
	logger := newLogger()

	// validate url
	if err := parseURL(*url); err != nil {
		os.Exit(1)
	}

	producerChan := make(chan string)
	consumerChan := make(chan string, noOfWorkers)

	client := notifier.NewHttpClient(logger, *url)
	notifier := notifier.NewNotifier(logger, client)
	// context
	ctx, cancel := context.WithCancel(context.Background())

	// read input
	go func(ctx context.Context, pChan chan string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scanner.Text()
		}

	}(ctx, produerChan)

	// send in to process
	wg := new(sync.Waitgroup)
	for i := 0; i < noOfWorkers; i++ {
		wg.add(1)
		go notifier.Process(ctx, wg)

	}
	//quit the process
	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	switch <-exit {
	case syscall.SIGINT, syscall.SIGTER:
		logger.Info("ctrl+d received, terminating")
	default:
		logger.Info("file read is complete")
	}
	signal.Stop(exit)
	cancel()
	wg.Wait()

}

func newLogger() *zap.Logger {
	cfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(colorable.NewColorableStdout()), zap.InfoLevel)
	return zap.New(core)
}

func parseURL(_url string) {

	parsedURL, err := url.Parse(url)
	if err != nil {
		return err
	}

}

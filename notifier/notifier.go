package notifier

import (
	"context"

	"go.uber.org/zap"
)

type Notifier struct {
	logger *zap.Logger
	client Client
	pChan  chan string
	cChan  chan string
}

func NewNotifier(logger *zap.Logger, client Client) *Notifier {
	return &Notifier{logger: logger, client: client}
}

func (n *Notifier) Process(ctx context.Context, wg *sync.Waitgroup) {
	defer wg.Done()
	for data := range <-n.cChan {
		n.client.Notify(data)
	}
}

func (n *Notifier) Start(ctx context.Context) {
	for {
		select {
		case data := <-n.pChan:
			n.cChan <- data
		case <-ctx.Done():
			close(n.pChan)
			close(n.cChan)
		}
	}
}

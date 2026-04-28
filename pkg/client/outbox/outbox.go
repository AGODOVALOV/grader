// Package outbox provides outbox functionality.
package outbox

import (
	"context"
	"sync"

	"github.com/AGODOVALOV/grader/pkg/client/user"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/queue/config"
	"github.com/streadway/amqp"
)

// Outbox represents the outbox.
type Outbox struct {
	userService *usecase.UserService
	cfg         *config.Config
	wg          *sync.WaitGroup
}

// NewOutbox creates a new outbox.
func NewOutbox(user *user.User, cfg *config.Config) *Outbox {
	return &Outbox{
		userService: user.Handler.Service,
		cfg:         cfg,
		wg:          &sync.WaitGroup{},
	}
}

// StartSending starts sending messages.
func (out *Outbox) StartSending(ctx context.Context) error {
	rConn, err := amqp.Dial(out.cfg.Broker.Rabbit.URL)
	if err != nil {
		logger.Z(ctx).Error(ctx, "rabbit connection", err.Error())
		return err
	}

	defer func(rConn *amqp.Connection) {
		err := rConn.Close()
		if err != nil {
			logger.Z(ctx).Error(ctx, "rabbit connection close", err.Error())
		}
	}(rConn)

	rCh, err := rConn.Channel()

	defer func(rCh *amqp.Channel) {
		err := rCh.Close()
		if err != nil {
			logger.Z(ctx).Error(ctx, "rabbit channel close", err.Error())
		}
	}(rCh)

	logger.Z(ctx).Info(ctx, "connection with rabbit is ok", "connection is ok")

	for _, cfg := range out.cfg.Messaging.Channels {
		out.wg.Go(func() {
			err = out.userService.ProduceMessages(ctx, rCh, &cfg)
			if err != nil {
				logger.Z(ctx).Error(ctx, "error producing messages", err.Error())
				return
			}
		})
	}
	out.wg.Wait()

	return nil
}

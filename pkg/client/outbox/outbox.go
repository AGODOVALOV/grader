// Package outbox provides outbox functionality.
package outbox

import (
	"context"
	"fmt"

	"github.com/AGODOVALOV/grader/pkg/client/user"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/queue/config"
	"github.com/streadway/amqp"
)

// Outbox represents the outbox.
type Outbox struct {
	user *user.User
	cfg  *config.Config
	rCh  *amqp.Channel
}

// NewOutbox creates a new outbox.
func NewOutbox(user *user.User, cfg *config.Config) *Outbox {
	return &Outbox{
		user: user,
		cfg:  cfg,
	}
}

// StartSending starts sending messages.
func (out *Outbox) StartSending(ctx context.Context) error {

	rCh, err := out.getChannel(ctx)
	if err != nil {
		return err
	}

	out.rCh = rCh

	logger.Z(ctx).Info(ctx, "connection with rabbit is ok", "connection is ok")

	return nil
}

func (out *Outbox) getChannel(ctx context.Context) (*amqp.Channel, error) {
	rConn, err := amqp.Dial(out.cfg.Broker.Rabbit.URL)
	if err != nil {
		logger.Z(ctx).Error(ctx, "rabbit connection", err.Error())
		return nil, err
	}
	defer func(rConn *amqp.Connection) {
		err := rConn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rConn)

	rCh, err := rConn.Channel()
	defer func(rCh *amqp.Channel) {
		err := rCh.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rCh)
	if err != nil {
		logger.Z(ctx).Error(ctx, "error open channel", err.Error())
		return nil, err
	}
	return rCh, nil
}

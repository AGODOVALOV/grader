package client

import (
	"bytes"
	"context"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/AGODOVALOV/grader/pkg/grader/client/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/token"
)

type Client struct {
	Client     *http.Client
	cfg        *config.Config
	tokenMaker token.Maker
}

func NewClient(cfg *config.Config, tokenMaker token.Maker) *Client {
	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConnections,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		TLSHandshakeTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Client{
		Client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		cfg:        cfg,
		tokenMaker: tokenMaker,
	}
}

func (c *Client) DoCallbackRequestWithRetry(ctx context.Context, payload []byte) error {
	const op = "grader.client.http.DoCallbackRequestWithRetry"
	var (
		dumpReq []byte
	)

	for attempt := 0; attempt <= c.cfg.Retry.MaxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			c.cfg.URL,
			bytes.NewReader(payload),
		)
		if err != nil {
			return err
		}

		jwtToken, _, err := c.tokenMaker.CreateToken(0, "grader")
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		dumpReq, _ = httputil.DumpRequest(req, false)

		resp, err := c.Client.Do(req)

		if err == nil && resp.StatusCode < 500 {
			logger.Z(ctx).Debug(ctx, op, "success", map[string]string{
				"attempt": strconv.Itoa(attempt + 1),
				"status":  resp.Status,
				"request": string(dumpReq),
			})
			return nil
		}

		if resp != nil {
			err = resp.Body.Close()
			if err != nil {
				return err
			}
		}

		if attempt < c.cfg.Retry.MaxAttempts {
			backoff := time.Duration(500<<uint(attempt)) * time.Millisecond
			jitter := time.Duration(rand.Intn(50)) * time.Millisecond

			logger.Z(ctx).Debug(ctx, op, "retrying request", map[string]string{
				"attempt": strconv.Itoa(attempt + 1),
				"request": string(dumpReq),
				"error":   getErrText(err),
			})

			select {
			case <-time.After(backoff + jitter):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

func getErrText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

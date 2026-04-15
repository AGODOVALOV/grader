package watcher

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/AGODOVALOV/grader/pkg/logger"
)

type Watcher struct{}

func NewWatcher() *Watcher {
	return &Watcher{}
}

func (w *Watcher) Watch(ctx context.Context, viper *viper.Viper, newConfig chan<- struct{}) {
	const op = "config.watcher.onconfigchange"

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		newConfig <- struct{}{}
		logger.Z(ctx).Debug(
			ctx,
			op,
			"Config file changed, reloading...",
			map[string]string{
				"file": e.Name,
			},
		)
	})
}

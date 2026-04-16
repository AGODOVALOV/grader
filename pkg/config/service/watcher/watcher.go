// Package watcher provides a watcher for config file changes
package watcher

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/AGODOVALOV/grader/pkg/logger"
)

// Watcher provides functionality to monitor configuration file changes and trigger reloading of the configuration.
type Watcher struct{}

// NewWatcher creates a new Watcher instance.
func NewWatcher() *Watcher {
	return &Watcher{}
}

// Watch starts watching for configuration file changes and notifies the provided channel when a change is detected.
func (*Watcher) Watch(ctx context.Context, v *viper.Viper, newConfig chan<- struct{}) {
	const op = "config.watcher.onconfigchange"

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
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

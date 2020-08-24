package config

import (
	"fmt"

	"github.com/micro/micro/v2/plugin"
)

var (
	defaultManager = plugin.NewManager()
)

// Plugins lists the config plugins
func Plugins() []plugin.Plugin {
	return defaultManager.Plugins()
}

// Register registers an config plugin
func Register(pl plugin.Plugin) error {
	if plugin.IsRegistered(pl) {
		return fmt.Errorf("%s registered globally", pl.String())
	}
	return defaultManager.Register(pl)
}

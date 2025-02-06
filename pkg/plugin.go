package pkg

import (
	"encoding/json"
	"github.com/Livepeer-Open-Pool/openpool-plugin/models"
)

// Plugin defines the interface for applying events to local state.
type Plugin interface {
	// Apply applies a batch of PoolEvent to the plugin's local state using the provided endpoint hash.
	// This ensures that only events fetched from that endpoint are allowed to update remote workers.
	Apply(events []models.PoolEvent, endpointHash string) error
	// GetEndpoints returns a list of endpoint URLs that should be polled.
	GetEndpoints() ([]string, error)
}

// PluginConstructor defines the standard constructor signature for plugins.
type PluginConstructor func(config json.RawMessage) (Plugin, error)

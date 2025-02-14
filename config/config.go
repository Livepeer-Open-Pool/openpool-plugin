package config

import (
	"encoding/json"
	"os"
)

// Config holds all the configuration parameters, serving as a unified config.
type Config struct {
	DataStorageFilePath    string                  `json:"DataStorageFilePath"`
	PluginPath             string                  `json:"PluginPath"`
	StoragePluginName      string                  `json:"StoragePluginName"`
	PoolCommissionRate     float64                 `json:"PoolCommissionRate"`
	Version                string                  `json:"Version"`
	Region                 string                  `json:"Region"`
	DataLoaderPluginConfig *DataLoaderPluginConfig `json:"DataLoaderPluginConfig,omitempty"`
	PayoutLoopConfig       *PayoutLoopConfig       `json:"PayoutLoopConfig,omitempty"`
	APIConfig              *APIConfig              `json:"APIConfig,omitempty"`
}

// LoadConfig reads a config file and unmarshals it into a unified Config struct.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type APIConfig struct {
	PluginName string `json:"PluginName"`
	ServerPort int    `json:"ServerPort,omitempty"`
}

type PayoutLoopConfig struct {
	PluginName               string `json:"PluginName"`
	RPCUrl                   string `json:"RPCUrl"`
	PrivateKeyStorePath      string `json:"PrivateKeyStorePath"`
	PrivateKeyPassphrasePath string `json:"PrivateKeyPassphrasePath"`
	PayoutFrequencySeconds   int    `json:"PayoutFrequencySeconds"`
	PayoutThreshold          string `json:"PayoutThreshold"`
}

// DataLoaderPluginConfig holds the plugin-specific settings for data loaders.
type DataLoaderPluginConfig struct {
	PluginName           string       `json:"PluginName"`
	FetchIntervalSeconds int          `json:"FetchIntervalSeconds"`
	DataSources          []DataSource `json:"Datasources"`
}

// DataSource defines the settings for each data source.
type DataSource struct {
	Endpoint string `json:"Endpoint"`
	NodeType string `json:"NodeType"`
}

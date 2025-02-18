package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"plugin"

	pool "github.com/Livepeer-Open-Pool/openpool-plugin"
	"github.com/Livepeer-Open-Pool/openpool-plugin/config"
)

// loadStoragePlugin loads the storage plugin and returns a StorageInterface.
func loadStoragePlugin(cfg *config.Config, pluginDir, storagePluginName string) (pool.StorageInterface, error) {
	if storagePluginName == "" {
		log.Fatal("Storage plugin name is not provided in the configuration.")
	}

	path := filepath.Join(pluginDir, storagePluginName)
	p, err := plugin.Open(path)
	if err != nil {
		log.Fatalf("Error loading storage plugin %s: %v", path, err)
	}

	symbol, err := p.Lookup("PluginInstance")
	if err != nil {
		log.Fatalf("Error finding symbol 'PluginInstance' in %s: %v", path, err)
	}

	storage, ok := symbol.(pool.StorageInterface)
	if !ok {
		log.Fatalf("Storage plugin does not implement StorageInterface")
	}

	storage.Init(cfg)
	return storage, nil
}

// loadGenericPlugin loads a plugin and returns an initialized PluginInterface.
func loadGenericPlugin(path string, cfg *config.Config, storage pool.StorageInterface) pool.PluginInterface {
	p, err := plugin.Open(path)
	if err != nil {
		log.Fatalf("Error loading plugin %s: %v", path, err)
	}

	symbol, err := p.Lookup("PluginInstance")
	if err != nil {
		log.Fatalf("Error finding symbol 'PluginInstance' in %s: %v", path, err)
	}

	instance, ok := symbol.(pool.PluginInterface)
	if !ok {
		log.Fatalf("Plugin %s does not implement PluginInterface", path)
	}

	instance.Init(*cfg, storage)
	return instance
}

func Run(configFileName string) {
	log.Printf("Starting Run with config file: %s\n", configFileName)

	// Load configuration from JSON file
	cfg, err := config.LoadConfig(configFileName)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Extract the plugin directory from config
	pluginDir := cfg.PluginPath
	if pluginDir == "" {
		log.Fatal("PluginPath is not provided in the configuration.")
	}

	// Load the storage plugin first
	fmt.Println("Loading storage plugin...")
	storage, err := loadStoragePlugin(cfg, pluginDir, cfg.StoragePluginName)
	if err != nil {
		log.Fatal("Failed to create storage plugin")
	}

	// Collect plugins dynamically from the config
	var pluginPaths []string

	if cfg.APIConfig != nil && cfg.APIConfig.PluginName != "" {
		log.Println("Loading Plugin from APIConfig")
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.APIConfig.PluginName))
	}

	if cfg.PayoutLoopConfig != nil && cfg.PayoutLoopConfig.PluginName != "" {
		log.Println("Loading Plugin from PayoutLoopConfig")
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.PayoutLoopConfig.PluginName))
	}

	if cfg.DataLoaderPluginConfig != nil && cfg.DataLoaderPluginConfig.PluginName != "" {
		log.Println("Loading Plugin from DataLoaderPluginConfig")
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.DataLoaderPluginConfig.PluginName))
	}

	var instances []pool.PluginInterface
	for _, path := range pluginPaths {
		instance := loadGenericPlugin(path, cfg, storage)
		instances = append(instances, instance)
	}
	log.Println("Plugins loaded successfully")

	// Start all plugins concurrently
	for idx, p := range instances {
		go func() {
			log.Println("Starting Plugin ", idx)
			p.Start()
		}()
	}
	log.Println("Plugins started successfully")

	select {} // Block forever
}

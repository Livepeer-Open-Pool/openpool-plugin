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
func loadStoragePlugin(cfg *config.Config, pluginDir, storagePluginName string) pool.StorageInterface {
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
	return storage
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
	//configFileName := flag.String("config", "/etc/pool/config.json", "Open Pool Configuration file to use")
	//flag.Parse()
	fmt.Printf("Using config file: %s\n", configFileName)

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
	storage := loadStoragePlugin(cfg, pluginDir, cfg.StoragePluginName)

	// Collect plugins dynamically from the config
	var pluginPaths []string

	if cfg.APIConfig != nil && cfg.APIConfig.PluginName != "" {
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.APIConfig.PluginName))
	}

	if cfg.PayoutLoopConfig != nil && cfg.PayoutLoopConfig.PluginName != "" {
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.PayoutLoopConfig.PluginName))
	}

	if cfg.DataLoaderPluginConfig != nil && cfg.DataLoaderPluginConfig.PluginName != "" {
		pluginPaths = append(pluginPaths, filepath.Join(pluginDir, cfg.DataLoaderPluginConfig.PluginName))
	}

	var instances []pool.PluginInterface
	for _, path := range pluginPaths {
		instance := loadGenericPlugin(path, cfg, storage)
		instances = append(instances, instance)
	}

	// Start all plugins concurrently
	for _, p := range instances {
		go p.Start()
	}

	fmt.Println("All plugins started.")
	select {} // Block forever
}

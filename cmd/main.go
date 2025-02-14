package main

import (
	"flag"
	"fmt"
	pool "github.com/Livepeer-Open-Pool/openpool-plugin"
	"github.com/Livepeer-Open-Pool/openpool-plugin/config"
	"log"
	"plugin"
)

// loadStoragePlugin loads the storage plugin and returns a StorageInterface.
func loadStoragePlugin(cfg *config.Config, path string) pool.StorageInterface {
	p, err := plugin.Open(path)
	if err != nil {
		log.Fatalf("Error loading storage plugin %s: %v", path, err)
	}

	symbol, err := p.Lookup("PluginInstance")
	if err != nil {
		log.Fatalf("Error finding symbol 'PluginInstance' in %s: %v", path, err)
	}

	// ✅ Correctly cast PluginInstance as StorageInterface
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

func main() {
	configFileName := flag.String("config", "/etc/pool/config.json", "Open Pool Configuration file to use")
	flag.Parse()
	fmt.Printf("Using config file: %s\n", *configFileName)

	// Load configuration from JSON file
	cfg, err := config.LoadConfig(*configFileName)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
		return
	}
	// Load the storage plugin first
	fmt.Println("Loading storage plugin...")
	//storage := loadStoragePlugin("../open-pool-basic-plugin/test-storage.so")
	storage := loadStoragePlugin(cfg, "../open-pool-basic-plugin/sqlite-storage.so")

	// Load other plugins
	pluginPaths := []string{
		"../open-pool-basic-plugin/dataloader.so",
		"../open-pool-basic-plugin/payoutloop.so",
		"../open-pool-basic-plugin/api.so",
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
	select {}
}

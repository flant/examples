package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Config fo our application
type Config struct {
	Welcome string `json:"welcome"`
	Name    string `json:"name"`
}

var (
	globalConfig *Config
)

// LoadConfig - load our config!
func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Unable to read configuration file %s", path)
	}

	config := new(Config)

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file %s", path)
	}

	return config, nil
}

// ConfigWatcher - watches config.json for changes
func ConfigWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				globalConfig, _ = LoadConfig("./configfiles/config.json")
				log.Println("config:", globalConfig)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./configfiles/config.json")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	log.Println("Start")
	globalConfig, _ = LoadConfig("./configfiles/config.json")
	go ConfigWatcher()
	for {
		log.Println("config:", globalConfig)
		time.Sleep(30 * time.Second)
	}
}

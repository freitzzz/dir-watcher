package main

import (
	"log"

	"github.com/freitzzz/dir-watcher/internal"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// 1. Parse rules.json
	rules, err := internal.Parse("rules.json")

	if err != nil {
		log.Fatal(err)
	}

	cache := internal.CacheMoveDirectories(rules.Move)

	// 2. Pre clean watching directories
	err = internal.AutoCleanDir(rules, cache)

	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	// 3. Start directories watcher.
	err = internal.Watch(rules, cache, watcher)

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

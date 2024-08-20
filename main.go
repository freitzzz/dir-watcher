package main

import (
	"log"

	"github.com/freitzzz/dir-watcher/internal"
	"github.com/fsnotify/fsnotify"
)

func main() {
	rules, err := internal.Parse("rules.json")

	if err != nil {
		log.Fatal(err)
	}

	c := internal.Cache(rules.Move)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	// Start listening for events.
	go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				filePath := event.Name

				if event.Has(fsnotify.Chmod) && !internal.ShouldIgnoreFile(filePath) {
					mp := c[internal.Ext(filePath)]

					if len(mp) == 0 {
						mp = internal.Expand(rules.Unknown)
					}

					err = internal.Move(filePath, mp)

					log.Println(err)
				}

				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, wp := range rules.Watch {
		ewp := internal.Expand(internal.GlobPath(wp))
		err = watcher.Add(ewp)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Watching %s directory\n", ewp)

		err = internal.CleanDir(ewp, c, internal.Expand(rules.Unknown))
	}

	<-make(chan struct{})
}

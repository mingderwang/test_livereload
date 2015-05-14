package main

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"log"
	"net/http"
)

// html includes the client JavaScript
const html = `<!doctype html>
<html>
<head>
  <title>Example</title>
</head>
<body>
  <p> hello reload when you touch /tmp/file </p>
  <script src="http://localhost:35729/livereload.js"></script>
</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// Add dir to watcher
	err = watcher.Add("/tmp")
	if err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr, _ := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	// Start goroutine that requests reload upon watcher event
	go func() {
		for {

			event := <-watcher.Events
			lr.Reload(event.Name)
		}
	}()

	go func() {
		fmt.Print("http://localhost:8080")

		// Start serving html
		http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(html))
		})
		//http.HandleFunc("/", handler)
		http.ListenAndServe(":8080", nil)

	}()
	<-done

}

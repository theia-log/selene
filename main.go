package main

import (
	"fmt"

	"github.com/theia-log/selene/watcher"
)

func main() {
	fmt.Println("test")

	wd := watcher.NewWatchDaemon()

	wd.Start()

	src := watcher.NewFileSource("./tmp/test")

	src.OnSourceEvent(func(file string, diff []byte) {
		fmt.Println(file, " -> ", string(diff))
	})
	wd.AddSource("ft", src)

	fmt.Println("started..")
	fmt.Scanf("L")
	wd.Stop()
}

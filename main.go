package main

import "sync"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	a, err := RegisterApplet("/home/charrington/go/bin/applet-demo")
	if err != nil {
		panic(err)
	}
	a.Start()
	wg.Wait()
}

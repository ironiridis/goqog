package main

import "log"

// CoreDispatcher is the general Applet message receiver. When fully implemented,
// it will handle requests from Applets by sending them to the appropriate
// handler functions.
func CoreDispatcher(a *Applet, m *AppletMessage) {
	log.Printf("CoreDispatcher: %#v, %#v\n", *a, *m)
	//TODO - future :(
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type AppletDispatcher func(*Applet, *AppletMessage)

// AppletMessage describes a message sent to or received from an applet over stdio.
type AppletMessage struct {
	T string
	D map[string]interface{}
}

// AppletState describes the current known or expected state of a defined applet.
type AppletState int

// AppletStateUnitialized is the zero value state of an Applet. AppletStateRegistered
// indicates an Applet that has been created via RegisterApplet but has not been Start()ed.
// AppletStateStarted indicates that the command has been been invoked, but goqog has not
// yet received the "ready" message from the applet. AppletStateRunning indicates that the
// Applet is running as expected. AppletStateFailed indicates that the Applet sent a
// failure message, but is still running. AppletStateCrashed indicates the Applet exited
// without being commanded to (or being killed). AppletStateStopped indicates the Applet
// was commanded to stop, or was killed.
const (
	AppletStateUninitialized AppletState = iota
	AppletStateRegistered
	AppletStateStarted
	AppletStateRunning
	AppletStateFailed
	AppletStateCrashed
	AppletStateStopped
)

// Applet describes a launchable program designed to communicate with the launcher.
type Applet struct {
	UUID     uuid.UUID
	invoke   string
	State    AppletState
	cmd      *exec.Cmd
	dispatch AppletDispatcher
	stdin    io.WriteCloser
	stdout   io.ReadCloser
	errchan  chan error
	stopchan chan struct{}
}

// RegisterApplet returns a new Applet struct that is ready to Start.
func RegisterApplet(invoke string) (a *Applet, err error) {
	a = &Applet{
		UUID:     uuid.New(),
		invoke:   invoke,
		cmd:      exec.Command(invoke),
		errchan:  make(chan error),
		stopchan: make(chan struct{}),
		dispatch: CoreDispatcher,
	}
	a.stdin, err = a.cmd.StdinPipe()
	if err != nil {
		return
	}
	a.stdout, err = a.cmd.StdoutPipe()
	if err != nil {
		return
	}
	a.cmd.Stderr = os.Stderr
	a.State = AppletStateRegistered
	return
}

func (a *Applet) readloop() {
	d := json.NewDecoder(a.stdout)
	for {
		select {
		case <-a.stopchan:
			return
		default:
			var m AppletMessage
			err := d.Decode(&m)
			if err != nil {
				a.errchan <- err
				return
			}
			a.dispatch(a, &m)
		}
	}
}

func (a *Applet) manager() {
	defer a.cmd.Wait()
	err := a.cmd.Start()
	if err != nil {
		a.errchan <- err
	}
	for {
		select {
		case <-a.stopchan:
			log.Printf("Applet stopping: %#v\n", *a)
			a.State = AppletStateStopped
			if a.stdin != nil {
				a.stdin.Close()
			}
			return
		case err := <-a.errchan:
			log.Printf("Applet error: %#v %#v\n", *a, err)
			close(a.stopchan)
		}
	}
}

func (a *Applet) Start() {
	go a.readloop()
	go a.manager()
}

func (a *Applet) Stop() {
	close(a.stopchan)

}

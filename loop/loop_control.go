package loop

import "log"

// loopController manages dynamic start/stop of background loops.
// Each loop runs in a goroutine and listens on its done channel.
// Closing the channel signals the loop to exit.
type LoopController struct {
	name     string
	doneChan chan struct{}
	active   bool
}

var (
	MonLoopCtl = &LoopController{name: "monitoring"}
	SecLoopCtl = &LoopController{name: "security"}
	RepLoopCtl = &LoopController{name: "reports"}
)

// start launches the loop goroutine if not already running.
func (lc *LoopController) Start(launchFunc func(chan struct{})) {
	if lc.active {
		log.Printf("[%s] Already running, skip", lc.name)
		return
	}
	lc.doneChan = make(chan struct{})
	lc.active = true
	go launchFunc(lc.doneChan)
	log.Printf("[%s] Loop started", lc.name)
}

// stop signals the loop to exit by closing the done channel.
func (lc *LoopController) Stop() {
	if !lc.active {
		log.Printf("[%s] Not running, skip stop", lc.name)
		return
	}
	close(lc.doneChan)
	lc.active = false
	log.Printf("[%s] Loop stopped", lc.name)
}

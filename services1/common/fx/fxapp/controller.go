package fxapp

type Controller struct {
	quitChannel chan bool
}

func (c *Controller) Quit() {
	c.quitChannel <- true
}

func newController() *Controller {
	return &Controller{
		quitChannel: make(chan bool, 1),
	}
}

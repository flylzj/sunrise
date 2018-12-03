package tk

import (
	"github.com/mattn/go-gtk/gtk"
	"time"
)

type Clock struct {
	TimeLabel *gtk.Label
}

func (c *Clock) showTime(){
	c.TimeLabel.SetText(time.Now().Format("15:04:05"))
}

func CreateClockWindow() *Clock {
	c := &Clock{}
	c.TimeLabel = gtk.NewLabel("")
	go func() {
		for {
			c.showTime()
			time.Sleep(time.Second)
		}
	}()
	return c
}
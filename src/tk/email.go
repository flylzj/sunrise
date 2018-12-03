package tk

import "github.com/mattn/go-gtk/gtk"

type Emailer struct {
	EmailFrame *gtk.Frame
	EmailBox   *gtk.VBox
	EmailEntry *gtk.Entry
}

func (e *Emailer) CreateMainFrame() {
	e.EmailBox.PackStart(e.EmailEntry, false, false, 0)
	e.EmailFrame.Add(e.EmailBox)
}

func CreateMailFrame() *Emailer {
	e := &Emailer{}
	e.EmailFrame = gtk.NewFrame("输入邮箱")

	e.EmailBox = gtk.NewVBox(false, 1)

	e.EmailEntry = gtk.NewEntry()

	e.CreateMainFrame()
	return e
	}

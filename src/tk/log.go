package tk

import (
	"github.com/mattn/go-gtk/gtk"
	"time"
)

type Log struct {
	LogFrame *gtk.Frame
	LogBox *gtk.VBox
	LogText *gtk.TextView
	LogScrolled *gtk.ScrolledWindow
	iter gtk.TextIter
}

func (l *Log) CreateLogFrame(){
	l.LogScrolled.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	l.LogScrolled.SetShadowType(gtk.SHADOW_IN)
	l.LogScrolled.Add(l.LogText)
	l.LogText.SetEditable(false)
	l.LogText.SetSizeRequest(800, 200)
	l.LogBox.PackStart(l.LogScrolled, false,false,0)
	l.LogFrame.Add(l.LogBox)
}

func (l *Log) Insert(text string){
	buffer := l.LogText.GetBuffer()

	buffer.GetEndIter(&l.iter)
	buffer.Insert(&l.iter, text)
}

func (l *Log) Update(){
	go func() {
		for {
			l.LogText.GetBuffer().SetText("")
			time.Sleep(time.Minute)
		}
	}()
}

func CreateLogFrame() *Log{
	l := &Log{}
	l.LogFrame = gtk.NewFrame("日志信息")
	l.LogScrolled = gtk.NewScrolledWindow(nil, nil)
	l.LogBox = gtk.NewVBox(false, 1)
	l.LogText = gtk.NewTextView()
	l.CreateLogFrame()
	l.Update()
	return l
}
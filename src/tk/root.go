package tk

import (
	"github.com/mattn/go-gtk/gtk"
	"os"
)
type Root struct {
	//主窗口
	Window *gtk.Window

	//窗口box
	MainBox *gtk.VBox

	MenuBar *gtk.MenuBar

	Notebook *gtk.Notebook

	MainFrame *gtk.Frame

	MainFrameBox *gtk.VBox

	MainFrameLog *Log

	Clocker *Clock

	Emailer *Emailer

	Exceler *Exceler

	OutputFrame *gtk.Frame

	OutputBox *gtk.VBox

	Spider *Spider

	OutputSpider *OutputSpider

	OutputLog *Log

	MessageChan chan [2]string
}

func (r *Root) CreateRootWindow(){
	r.Window.SetTitle("日上会员")
	r.Window.SetSizeRequest(800, 600)

	r.MainFrameBox.PackStart(r.Clocker.TimeLabel, false, false, 0)
	r.MainFrameBox.PackStart(r.Emailer.EmailFrame, false, false, 0)
	r.MainFrameBox.PackStart(r.Exceler.ExcelFrame, false, false, 0)
	r.MainFrameBox.PackStart(r.Spider.SpiderFrame, false, false,0)
	r.MainFrameBox.PackStart(r.MainFrameLog.LogFrame, false, false,0)

	r.MainFrame.Add(r.MainFrameBox)

	r.OutputBox.PackStart(r.OutputSpider.OutputFrame,false,false,0)
	r.OutputBox.PackStart(r.OutputLog.LogFrame, false,false,0)

	r.OutputFrame.Add(r.OutputBox)

	r.Notebook.AppendPage(r.MainFrame, gtk.NewLabel("主页"))
	r.Notebook.AppendPage(r.OutputFrame, gtk.NewLabel("导出"))

	r.MainBox.PackStart(r.MenuBar, false, false, 0)
	r.MainBox.PackStart(r.Notebook, false, false, 0)

	r.Window.Add(r.MainBox)
}

func (r *Root) InsertMessage(){
	for {
		message := <- r.MessageChan
		switch message[0] {
		case "main":
			r.MainFrameLog.Insert(message[1]+"\n")
		case "output":
			r.OutputLog.Insert(message[1]+"\n")
		}
	}
}


func MakeRootWindow()*Root{
	gtk.Init(&os.Args)
	r := &Root{}
	r.Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	r.MainBox = gtk.NewVBox(false, 1)
	r.MenuBar = CreateMenu()
	r.Notebook = CreateNotebook()
	r.MainFrame = gtk.NewFrame("")
	r.MainFrameBox = gtk.NewVBox(false, 1)
	r.Clocker = CreateClockWindow()
	r.Emailer = CreateMailFrame()
	r.Exceler = CreateExcelFrame()
	r.Spider = CreateSpiderFrame()
	r.MainFrameLog = CreateLogFrame()
	r.OutputFrame = gtk.NewFrame("")
	r.OutputBox = gtk.NewVBox(false, 1)
	r.OutputSpider = CreateOutputFrame()
	r.OutputLog = CreateLogFrame()
	r.CreateRootWindow()
	r.MessageChan = make(chan [2]string, 99)
	go r.InsertMessage()
	return r
}

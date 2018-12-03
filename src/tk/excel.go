package tk

import (
	"github.com/mattn/go-gtk/gtk"
)

type Exceler struct {
	ExcelFrame      *gtk.Frame
	ExcelBox        *gtk.VBox
	ExcelEntry      *gtk.Entry
	ChoseFileButton *gtk.Button
}

func (e *Exceler) CreateExcelFrame(){
	e.ExcelBox.PackStart(e.ExcelEntry, false, false, 0)
	e.ExcelBox.PackStart(e.ChoseFileButton, false, false, 0)
	e.ExcelFrame.Add(e.ExcelBox)
}

func CreateExcelFrame() *Exceler {
	e := &Exceler{}

	e.ExcelFrame = gtk.NewFrame("选择文件")

	e.ExcelBox = gtk.NewVBox(false, 1)

	e.ExcelEntry = gtk.NewEntry()

	e.ChoseFileButton = gtk.NewButtonWithLabel("选择文件")
	e.ChoseFileButton.Clicked(func() {
			filechooserdialog := gtk.NewFileChooserDialog(
				"Choose File...",
				e.ChoseFileButton.GetTopLevelAsWindow(),
				gtk.FILE_CHOOSER_ACTION_OPEN,
				gtk.STOCK_OK,
				gtk.RESPONSE_ACCEPT)
			filter := gtk.NewFileFilter()
			filter.AddPattern("*.xlsx")
			filechooserdialog.AddFilter(filter)
			filechooserdialog.Response(func() {
				e.ExcelEntry.SetText(filechooserdialog.GetFilename())
				filechooserdialog.Destroy()
			})
			filechooserdialog.Run()
	})
	e.CreateExcelFrame()
	return e
}
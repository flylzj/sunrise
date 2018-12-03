package tk

import "github.com/mattn/go-gtk/gtk"

func CreateNotebook() *gtk.Notebook{
	notebook := gtk.NewNotebook()
	return notebook
}

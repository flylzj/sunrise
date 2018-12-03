package tk

import "github.com/mattn/go-gtk/gtk"

func CreateMenu() *gtk.MenuBar{
	menuBar := gtk.NewMenuBar()

	cascadeMenu := gtk.NewMenuItemWithMnemonic("_选项")
	menuBar.Append(cascadeMenu)
	submenu := gtk.NewMenu()
	cascadeMenu.SetSubmenu(submenu)

	menuItem := gtk.NewMenuItemWithMnemonic("退_出")
	menuItem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuItem)

	return menuBar
}

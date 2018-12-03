package tk

import (
	"github.com/mattn/go-gtk/gtk"
)

type Spider struct {
	SpiderFrame *gtk.Frame
	SpiderBox   *gtk.VBox

	HourSpinButton   *gtk.SpinButton
	MinuteSpinButton *gtk.SpinButton
	TimeLabel        *gtk.Label
	HourLabel        *gtk.Label
	MinuteLabel      *gtk.Label

	StartButton  *gtk.Button
}

func (s *Spider) CreateSpiderFrame() {
	s.SpiderBox.PackStart(s.TimeLabel, false, false, 0)
	s.SpiderBox.PackStart(s.HourLabel, false, false, 0)
	s.SpiderBox.PackStart(s.HourSpinButton, false, false, 0)
	s.SpiderBox.PackStart(s.MinuteLabel, false, false,0)
	s.SpiderBox.PackStart(s.MinuteSpinButton, false, false, 0)

	s.SpiderBox.PackStart(s.StartButton, false, false, 0)

	s.SpiderFrame.Add(s.SpiderBox)
}

func CreateSpiderFrame() *Spider {
	s := &Spider{}
	s.SpiderFrame = gtk.NewFrame("操作控制")

	s.SpiderBox = gtk.NewVBox(false, 0)

    //爬取周期控制
	s.HourSpinButton = gtk.NewSpinButtonWithRange(0, 23, 1)
	s.MinuteSpinButton = gtk.NewSpinButtonWithRange(2, 59, 1)
	s.TimeLabel = gtk.NewLabel("请输入爬取周期:")
	s.HourLabel = gtk.NewLabel("小时:")
	s.MinuteLabel = gtk.NewLabel("分钟:")


	//操作按钮
	s.StartButton = gtk.NewButtonWithLabel("开始爬取")

	s.CreateSpiderFrame()
	return s
}

type OutputSpider struct {
	OutputFrame  *gtk.Frame
	OutputBox    *gtk.VBox
	OutputButton *gtk.Button
}

func (os *OutputSpider) CreateOutputFrame() {
	os.OutputBox.PackStart(os.OutputButton, false, false, 0)
	os.OutputFrame.Add(os.OutputBox)
}

func CreateOutputFrame() *OutputSpider {
	os := &OutputSpider{}
	os.OutputFrame = gtk.NewFrame("导出全部商品:")

	os.OutputBox = gtk.NewVBox(false, 1)
	os.OutputButton = gtk.NewButtonWithLabel("导出")

	os.CreateOutputFrame()
	return os
}

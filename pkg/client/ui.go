package client

import (
	"fmt"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type OnSubmitCallback func(message string)

type Ui struct {
	MsgList MsgList

	input *widgets.Paragraph
	list  *widgets.List

	onSubmit OnSubmitCallback

	Lock sync.Mutex
}

func NewUi(callback OnSubmitCallback) *Ui {
	return &Ui{
		MsgList:  MsgList{nil, nil, 0},
		onSubmit: callback,
	}
}

func (u *Ui) Run() {
	errTer := ui.Init()
	if errTer != nil {
		panic(errTer)
	}
	defer ui.Close()

	u.list = widgets.NewList()
	u.list.Title = "Messages"
	u.list.Rows = u.MsgList.MessageArr()
	u.list.WrapText = false
	u.list.SetRect(0, 0, 50, 8)

	u.input = widgets.NewParagraph()
	u.input.Text = ">"
	u.input.SetRect(0, 9, 50, 13)

	u.render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Enter>":
			u.onEnter()
			break
		case "<Backspace>":
			len := len(u.input.Text)
			if len > 1 {
				u.input.Text = u.input.Text[:len-1]
			}
			break
		case "<Space>":
			u.input.Text += " "
			break
		default:
			u.input.Text += e.ID
			break
		}

		u.render()
	}
}

func (u *Ui) render() {
	ui.Render(u.list, u.input)
}

func (u *Ui) onEnter() {
	u.onSubmit(u.input.Text[1:len(u.input.Text)])
	u.input.Text = ">"
}

func (u *Ui) AddMessage(message string) {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	if u.list == nil {
		fmt.Println("null")
		return
	}
	u.MsgList.Append(message)

	u.list.Rows = u.MsgList.MessageArr()
	u.list.ScrollBottom()

	u.render()
}

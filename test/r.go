package test

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("自定义对话框示例")

	button := widget.NewButton("打开对话框", func() {
		content := container.NewVBox(widget.NewLabel("这是一个自定义对话框"), widget.NewLabel("这是一个自定义对话框"), widget.NewLabel("这是一个自定义对话框"), widget.NewLabel("这是一个自定义对话框"), widget.NewLabel("这是一个自定义对话框"), widget.NewLabel("这是一个自定义对话框"))
		dialog.ShowCustom("标题", "关闭", content, myWindow)
	})

	myWindow.SetContent(container.NewVBox(
		button,
	))

	myWindow.ShowAndRun()
}

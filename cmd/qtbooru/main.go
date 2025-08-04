package main

import (
	"fmt"
	"os"

	q "github.com/mappu/miqt/qt6"
)

const (
	itemWidth int = 180
	itemHeight int = 194
)

func main() {
	ui()
}

func ui() {
	q.NewQApplication(os.Args)

	window := q.NewQMainWindow2()

	scrollarea := q.NewQScrollArea(window.QWidget)
	content := q.NewQWidget(scrollarea.QWidget)
	layout := q.NewQGridLayout2()
	content.SetLayout(layout.Layout())

	width := window.Width() / itemWidth
	items := make([]*q.QWidget, 0, 20)
	for i := range 20 {
		item := q.NewQPushButton5(fmt.Sprint(i), content)
		item.Show()
		item.SetFixedWidth(itemWidth)
		item.SetFixedHeight(itemHeight)
		items = append(items, item.QWidget)
		layout.AddWidget2(item.QWidget, i / width, i % width)
	}

	content.Show()
	window.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		relayout(layout, items, (max(itemWidth, event.Size().Width() - 20)) / itemWidth)
	})
	window.SetMinimumSize2(itemWidth, itemHeight)
	window.SetCentralWidget(scrollarea.QWidget)
	window.Show()
	q.QApplication_Exec()
}

func relayout(layout *q.QGridLayout, items []*q.QWidget, width int) {
	parent := items[0].ParentWidget()
	for _, item := range items {
		layout.RemoveWidget(item)
	}
	rows := 0
	for i, item := range items {
		rows = i / width
		layout.AddWidget2(item, rows, i % width)
	}
	rows++
	w := parent.ParentWidget().Width()
	h := rows * itemHeight
	h += rows*20
	parent.SetFixedSize2(w, h)
}


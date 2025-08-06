package main

import (
	"log"
	"net/http"
	"os"
	"qtbooru/pkg/api"
	"qtbooru/pkg/api/post"

	"github.com/joho/godotenv"
	q "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

const (
	itemWidth int = 200
	itemHeight int = 200

	Agent = "QtBooru/indev_v0 (created by readf0x)"
)
var client = &http.Client{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	tags := os.Args[1:]
	posts := *(&api.RequestBuilder{
		Site: api.E926,
		Params: &[]string{"limit=20"},
		Tags: &tags,
		Key: os.Getenv("API_USER"),
		Agent: os.Getenv("API_KEY"),
	}).Process(client)

	q.NewQApplication([]string{})

	window := q.NewQMainWindow2()

	scrollarea := q.NewQScrollArea(window.QWidget)
	scrollarea.SetWidgetResizable(true)
	scrollarea.SetHorizontalScrollBarPolicy(q.ScrollBarAlwaysOff)
	content := q.NewQWidget(scrollarea.QWidget)
	layout := q.NewQGridLayout2()
	content.SetLayout(layout.Layout())
	scrollarea.SetWidget(content)

	width := window.Width() / itemWidth
	items := make([]*q.QWidget, 0, len(posts))
	for i, p := range posts {
		item := q.NewQLabel5(p.Description, content)
		go getAsync(&p.Preview, item)
		item.SetFixedWidth(itemWidth)
		item.SetFixedHeight(itemHeight)
		item.Show()
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

func getAsync(f *post.File, label *q.QLabel) {
	b, err := f.Get(client)
	if err != nil { return }
	image := q.NewQPixmap()
	image.LoadFromDataWithData(*b)
	mainthread.Wait(func() {
		label.SetPixmap(image)
	})
}

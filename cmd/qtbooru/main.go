package main

import (
	"log"
	"net/http"
	"os"
	"qtbooru/pkg/api"

	"github.com/joho/godotenv"
	q "github.com/mappu/miqt/qt6"
)

const (
	itemWidth int = 200
	itemHeight int = 200
)
var client = &http.Client{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	tags := os.Args[1:]
	req, err := api.NewRequest(
		api.E926,
		&[]string{"limit=20"},
		&tags,
		os.Getenv("API_USER"),
		os.Getenv("API_KEY"),
	)

	posts := *api.Process(client, req)

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
	for i, post := range posts {
		// fmt.Println(post.Preview.URL)
		item := q.NewQLabel5(post.Description, content)

		b, err := post.Preview.Get(client)
		if err != nil { log.Fatal(err) }
		image := q.NewQPixmap2(post.Preview.Width, post.Preview.Height)
		image.LoadFromDataWithData(*b)

		item.SetPixmap(image)

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


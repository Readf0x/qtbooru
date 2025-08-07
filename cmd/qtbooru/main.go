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
	req := &api.RequestBuilder{
		Site: api.E926,
		Params: &[]string{"limit=20"},
		Tags: &tags,
		User: os.Getenv("API_USER"),
		Key: os.Getenv("API_KEY"),
		Agent: Agent,
	}
	p, err := req.Process(client)
	if err != nil {
		log.Fatal(err)
	}
	posts := *p

	q.NewQApplication([]string{})

	window := q.NewQMainWindow2()
	stack := q.NewQStackedWidget(window.QWidget)
	window.SetCentralWidget(stack.QWidget)

	itemList := q.NewQScrollArea(window.QWidget)
	itemList.SetWidgetResizable(true)
	itemList.SetHorizontalScrollBarPolicy(q.ScrollBarAlwaysOff)
	content := q.NewQWidget(itemList.QWidget)
	layout := q.NewQGridLayout2()
	itemList.SetWidget(content)
	stack.AddWidget(itemList.QWidget)

	imageView := q.NewQLabel(window.QWidget)
	imageView.SetVisible(false)
	imageView.SetScaledContents(true)
	imageView.SetSizePolicy2(q.QSizePolicy__Ignored, q.QSizePolicy__Ignored)
	stack.AddWidget(imageView.QWidget)

	width := window.Width() / itemWidth
	var items *[]*q.QWidget = nil
	if len(posts) == 0 {
		hlayout := q.NewQHBoxLayout2()
		msg := q.NewQLabel3("No Results.")
		msg.SetAlignment(q.AlignCenter)
		hlayout.AddWidget(msg.QWidget)
		hlayout.SetAlignment(msg.QWidget, q.AlignCenter)
		content.SetLayout(hlayout.Layout())
	} else {
		content.SetLayout(layout.Layout())
		it := make([]*q.QWidget, 0, len(posts))
		items = &it
		for i, p := range posts {
			item := q.NewQLabel5(p.Description, content)
			go getAsync(&p.Preview, item)
			item.SetFixedWidth(itemWidth)
			item.SetFixedHeight(itemHeight)
			item.OnMousePressEvent(func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent){
				if ev.Button() == q.LeftButton {
					var f *post.File
					if p.Sample.URL != "" {
						f = &p.Sample.File
					} else {
						f = &p.File.File
					}
					go getAsync(f, imageView)
					imageView.SetVisible(true)
					stack.SetCurrentIndex(1)
				}
				super(ev)
			})
			it = append(it, item.QWidget)
			layout.AddWidget2(item.QWidget, i / width, i % width)
		}
	}

	itemList.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		relayout(layout, items, (max(itemWidth, event.Size().Width() - 20)) / itemWidth)
	})
	window.SetMinimumSize2(itemWidth, itemHeight)
	window.OnKeyPressEvent(func(super func(event *q.QKeyEvent), event *q.QKeyEvent){
		if event.Key() == int(q.Key_Escape) {
			stack.SetCurrentIndex(0)
		}
	})
	window.Show()
	q.QApplication_Exec()
}

func relayout(layout *q.QGridLayout, items *[]*q.QWidget, width int) {
	if items != nil {
		i := *items
		parent := i[0].ParentWidget()
		for _, item := range i {
			layout.RemoveWidget(item)
		}
		rows := 0
		for i, item := range i {
			rows = i / width
			layout.AddWidget2(item, rows, i % width)
		}
		rows++
		w := parent.ParentWidget().Width()
		h := rows * itemHeight
		h += rows*20
		parent.SetFixedSize2(w, h)
	}
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

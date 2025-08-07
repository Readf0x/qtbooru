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

	posts := request(os.Args[1:])

	q.NewQApplication([]string{})

	window := q.NewQMainWindow2()
	stack := q.NewQStackedWidget(window.QWidget)
	window.SetCentralWidget(stack.QWidget)

	mainArea := q.NewQWidget(window.QWidget)
	mainLayout := q.NewQVBoxLayout(mainArea)
	mainArea.SetLayout(mainLayout.Layout())
	search := q.NewQLineEdit2()
	mainLayout.AddWidget(search.QWidget)
	itemList := q.NewQScrollArea2()
	itemList.SetWidgetResizable(true)
	itemList.SetHorizontalScrollBarPolicy(q.ScrollBarAlwaysOff)
	mainLayout.AddWidget(itemList.QWidget)
	listContent := q.NewQWidget(itemList.QWidget)
	layout := q.NewQGridLayout2()
	listContent.SetLayout(layout.Layout())
	itemList.SetWidget(listContent)
	stack.AddWidget(mainArea)

	imageView := q.NewQLabel(window.QWidget)
	imageView.SetVisible(false)
	imageView.SetScaledContents(true)
	imageView.SetSizePolicy2(q.QSizePolicy__Ignored, q.QSizePolicy__Ignored)
	stack.AddWidget(imageView.QWidget)

	width := window.Width() / itemWidth
	var items []*q.QWidget = nil
	if len(posts) == 0 {
		msg := q.NewQLabel3("No Results.")
		msg.SetAlignment(q.AlignCenter)
		listContent.Layout().AddWidget(msg.QWidget)
	} else {
		items = make([]*q.QWidget, 0, len(posts))
		posts.addToGrid(&items, width, layout, imageView, stack)
	}

	itemList.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		if len(items) != 0 {
			relayout(layout, items, max(itemWidth, event.Size().Width() - 20) / itemWidth)
		}
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

func relayout(layout *q.QGridLayout, items []*q.QWidget, width int) {
	if items != nil {
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

func request(tags []string) posts {
	req := &api.RequestBuilder{
		Site: api.E926,
		Params: &[]string{"limit=4"},
		Tags: &tags,
		User: os.Getenv("API_USER"),
		Key: os.Getenv("API_KEY"),
		Agent: Agent,
	}
	p, err := req.Process(client)
	if err != nil {
		log.Fatal(err)
	}
	return *p
}

type posts []*post.Post

func (P posts) addToGrid(items *[]*q.QWidget, width int, grid *q.QGridLayout, imageView *q.QLabel, stack *q.QStackedWidget) {
	for i, p := range P {
		item := q.NewQLabel3(p.Description)
		go getAsync(&p.Preview, item)
		item.SetFixedWidth(itemWidth)
		item.SetFixedHeight(itemHeight)
		item.OnMousePressEvent(itemClick(p, imageView, stack))
		*items = append(*items, item.QWidget)
		grid.AddWidget2(item.QWidget, i / width, i % width)
	}
}

func itemClick(p *post.Post, imageView *q.QLabel, stack *q.QStackedWidget) func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent) {
	return func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent){
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
		}
}

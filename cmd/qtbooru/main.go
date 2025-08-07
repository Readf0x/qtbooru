package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"qtbooru/pkg/api"
	"qtbooru/pkg/api/post"
	"strings"

	"github.com/joho/godotenv"
	q "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

const (
	itemWidth = 200
	itemHeight = 200
	pageSize = 40

	Agent = "QtBooru/indev_v0 (created by readf0x)"
)
var (
	currentPage = 1
	isLoading = false
	endOfPosts = false

	client = &http.Client{}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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
	scrollBar := itemList.VerticalScrollBar()
	mainLayout.AddWidget(itemList.QWidget)
	listContent := q.NewQWidget(itemList.QWidget)
	layout := q.NewQGridLayout2()
	listContent.SetLayout(layout.Layout())
	itemList.SetWidget(listContent)
	stack.AddWidget(mainArea)

	imageView := q.NewQGraphicsView(window.QWidget)
	imageView.SetRenderHint(q.QPainter__SmoothPixmapTransform)
	imageView.SetTransformationAnchor(q.QGraphicsView__AnchorUnderMouse)
	imageView.SetResizeAnchor(q.QGraphicsView__AnchorUnderMouse)
	imageView.SetDragMode(q.QGraphicsView__ScrollHandDrag)
	imageView.SetHorizontalScrollBarPolicy(q.ScrollBarAlwaysOff)
	imageView.SetVerticalScrollBarPolicy(q.ScrollBarAlwaysOff)
	imageView.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		if scene := imageView.Scene(); scene != nil {
			if !scene.SceneRect().IsNull() {
				imageView.FitInView3(scene.SceneRect(), q.KeepAspectRatio)
			}
		}
		super(event)
	})
	scene := q.NewQGraphicsScene()
	imageView.SetScene(scene)
	stack.AddWidget(imageView.QWidget)

	items := make([]*q.QWidget, 0)

	tags := []string{}
	if len(os.Args) > 1 {
		tags = os.Args[1:]
		search.SetText(strings.Join(tags, " "))
	}

	update(tags, &items, listContent, layout, imageView, stack)

	search.OnReturnPressed(func(){
		tags = strings.Split(search.Text(), " ")
		for _, item := range items {
			layout.RemoveWidget(item)
			item.DeleteLater()
		}
		items = make([]*q.QWidget, 0)
		currentPage = 1
		scrollBar.SetValue(0)
		update(tags, &items, listContent, layout, imageView, stack)
	})
	itemList.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		if len(items) != 0 {
			relayout(listContent, layout, items, max(itemWidth, event.Size().Width() - 20) / itemWidth)
		}
	})
	window.SetMinimumSize2(itemWidth, itemHeight)
	window.OnKeyPressEvent(func(super func(event *q.QKeyEvent), event *q.QKeyEvent){
		if event.Key() == int(q.Key_Escape) {
			stack.SetCurrentIndex(0)
		}
	})
	scrollBar.OnValueChanged(func(value int){
		if isLoading || endOfPosts {
			return
		}
		if value >= scrollBar.Maximum() {
			update(tags, &items, listContent, layout, imageView, stack)
		}
	})
	window.Show()
	q.QApplication_Exec()
}

func relayout(parent *q.QWidget, layout *q.QGridLayout, items []*q.QWidget, width int) {
	if len(items) > 0 {
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
		h += rows * pageSize
		parent.SetFixedSize2(w, h)
	}
	layout.Activate()
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

func viewFile(f *post.File, view *q.QGraphicsView) {
	b, err := f.Get(client)
	if err != nil { return }
	image := q.NewQPixmap()
	image.LoadFromDataWithData(*b)
	mainthread.Wait(func() {
		scene := view.Scene()
		scene.Clear()
		item := scene.AddPixmap(image)
		item.SetTransformationMode(q.FastTransformation)
		scene.SetSceneRect2(0, 0, float64(image.Width()), float64(image.Height()))
		view.FitInView3(scene.SceneRect(), q.KeepAspectRatio)
	})
}

func request(tags []string) posts {
	req := &api.RequestBuilder{
		Site: api.E926,
		Params: &[]string{fmt.Sprintf("limit=%d", pageSize), fmt.Sprintf("page=%d", currentPage)},
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

func (P posts) addToGrid(items *[]*q.QWidget, grid *q.QGridLayout, imageView *q.QGraphicsView, stack *q.QStackedWidget) {
	for _, p := range P {
		item := q.NewQLabel3(p.Description)
		go getAsync(&p.Preview, item)
		item.SetFixedWidth(itemWidth)
		item.SetFixedHeight(itemHeight)
		item.OnMousePressEvent(itemClick(p, imageView, stack))
		*items = append(*items, item.QWidget)
		grid.AddWidget(item.QWidget)
	}
}

func itemClick(p *post.Post, imageView *q.QGraphicsView, stack *q.QStackedWidget) func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent) {
	return func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent){
			if ev.Button() == q.LeftButton {
				var f *post.File
				if p.Sample.URL != "" {
					f = &p.Sample.File
				} else {
					f = &p.File.File
				}
				go viewFile(f, imageView)
				stack.SetCurrentIndex(1)
			}
			super(ev)
		}
}

func update(tags []string, items *[]*q.QWidget, listContent *q.QWidget, layout *q.QGridLayout, imageView *q.QGraphicsView, stack *q.QStackedWidget) {
	isLoading = true
	posts := request(tags)

	if len(posts) < pageSize { endOfPosts = true }

	if len(*items) == 0 && len(posts) == 0 {
		msg := q.NewQLabel3("No Results.")
		msg.SetAlignment(q.AlignCenter)
		listContent.Layout().AddWidget(msg.QWidget)
	} else {
		posts.addToGrid(items, layout, imageView, stack)
		isLoading = false
		currentPage++
		relayout(listContent, layout, *items, max(itemWidth, listContent.Window().Width() - 20) / itemWidth)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"qtbooru/config"
	"qtbooru/pkg/api"
	"qtbooru/pkg/api/post"
	"strings"

	q "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
	m "github.com/mappu/miqt/qt6/multimedia"
)

const (
	itemWidth  = 200
	itemHeight = 200
	pageSize   = 40

	Agent = "QtBooru/v1.0 (created by readf0x)"
)

var (
	currentPage = 1
	isLoading   = false
	endOfPosts  = false
	client      = &http.Client{}
	booru       = api.E926
	initialTags []string
)

func main() {
	conf := loadEnv()

	processArgs(os.Args)

	q.NewQApplication([]string{})

	window := q.NewQMainWindow2()
	stack := q.NewQStackedWidget(window.QWidget)
	window.SetCentralWidget(stack.QWidget)

	mainArea := q.NewQWidget(window.QWidget)
	mainLayout := q.NewQVBoxLayout(mainArea)
	mainArea.SetLayout(mainLayout.Layout())
	search := q.NewQLineEdit2()
	search.SetStyleSheet("font-size: 12pt;")
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

	player := m.NewQMediaPlayer()
	videoView := m.NewQVideoWidget(window.QWidget)
	player.SetVideoOutput(videoView.QObject)
	player.SetAudioOutput(m.NewQAudioOutput2(m.QMediaDevices_DefaultAudioOutput()))
	stack.AddWidget(videoView.QWidget)
	stack.OnCurrentChanged(func(p int) {
		if p == 0 {
			player.Stop()
		}
	})

	items := make([]*q.QWidget, 0)

	tags := initialTags
	search.SetText(strings.Join(tags, " "))

	update(tags, conf, &items, listContent, layout, imageView, player, stack)

	search.OnReturnPressed(func() {
		tags = strings.Split(search.Text(), " ")
		for _, item := range items {
			layout.RemoveWidget(item)
			item.DeleteLater()
		}
		items = make([]*q.QWidget, 0)
		currentPage = 1
		scrollBar.SetValue(0)
		update(tags, conf, &items, listContent, layout, imageView, player, stack)
	})
	itemList.OnResizeEvent(func(super func(event *q.QResizeEvent), event *q.QResizeEvent) {
		if len(items) != 0 {
			relayout(listContent, layout, items, max(itemWidth, event.Size().Width()-20)/itemWidth)
		}
	})
	window.SetMinimumSize2(itemWidth, itemHeight)
	window.OnKeyPressEvent(func(super func(event *q.QKeyEvent), event *q.QKeyEvent) {
		if event.Key() == int(q.Key_Escape) {
			stack.SetCurrentIndex(0)
		}
	})
	scrollBar.OnValueChanged(func(value int) {
		if isLoading || endOfPosts {
			return
		}
		if value >= scrollBar.Maximum() {
			update(tags, conf, &items, listContent, layout, imageView, player, stack)
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
			layout.AddWidget2(item, rows, i%width)
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
	if err != nil {
		return
	}
	image := q.NewQPixmap()
	image.LoadFromDataWithData(*b)
	mainthread.Wait(func() {
		label.SetPixmap(image)
	})
}

func viewFile(f *post.File, view *q.QGraphicsView) {
	b, err := f.Get(client)
	if err != nil {
		return
	}
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

func viewVideo(f *post.File, player *m.QMediaPlayer) {
	player.SetSource(q.NewQUrl3(strings.Replace(string(f.URL), "localhost", "10.1.11.104", 1)))
	player.Play()
}

func request(tags []string, conf *config.ApiConfig) posts {
	req := &api.RequestBuilder{
		Site:   booru,
		Params: &[]string{fmt.Sprintf("limit=%d", pageSize), fmt.Sprintf("page=%d", currentPage)},
		Tags:   &tags,
		User:   conf.Username,
		Key:    conf.Key,
		Agent:  Agent,
	}
	p, err := req.Process(client)
	if err != nil {
		log.Fatal(err)
	}
	return *p
}

type posts []*post.Post

func (P posts) addToGrid(items *[]*q.QWidget, grid *q.QGridLayout, imageView *q.QGraphicsView, player *m.QMediaPlayer, stack *q.QStackedWidget) {
	for _, p := range P {
		item := q.NewQLabel3(p.Description)
		go getAsync(&p.Preview, item)
		item.SetFixedWidth(itemWidth)
		item.SetFixedHeight(itemHeight)
		item.OnMousePressEvent(itemClick(p, imageView, player, stack))
		*items = append(*items, item.QWidget)
		grid.AddWidget(item.QWidget)
	}
}

func itemClick(p *post.Post, imageView *q.QGraphicsView, player *m.QMediaPlayer, stack *q.QStackedWidget) func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent) {
	return func(super func(ev *q.QMouseEvent), ev *q.QMouseEvent) {
		if ev.Button() == q.LeftButton {
			var f *post.File
			if p.Sample.URL != "" {
				f = &p.Sample.File
			} else {
				f = &p.File.File
			}
			if p.File.Type == post.WEBM {
				viewVideo(&p.File.File, player)
				stack.SetCurrentIndex(2)
			} else {
				go viewFile(f, imageView)
				stack.SetCurrentIndex(1)
			}
		}
		super(ev)
	}
}

func update(tags []string, config *config.ApiConfig, items *[]*q.QWidget, listContent *q.QWidget, layout *q.QGridLayout, imageView *q.QGraphicsView, player *m.QMediaPlayer, stack *q.QStackedWidget) {
	isLoading = true
	posts := request(tags, config)

	if len(posts) < pageSize {
		endOfPosts = true
	}

	if len(*items) == 0 && len(posts) == 0 {
		msg := q.NewQLabel3("No Results.")
		msg.SetAlignment(q.AlignCenter)
		listContent.Layout().AddWidget(msg.QWidget)
	} else {
		posts.addToGrid(items, layout, imageView, player, stack)
		isLoading = false
		currentPage++
		relayout(listContent, layout, *items, max(itemWidth, listContent.Window().Width()-20)/itemWidth)
	}
}

func processArgs(args []string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-b":
			fallthrough
		case "--booru":
			i++
			if i >= len(args) {
				log.Fatal("Missing booru after --booru")
			}
			switch args[i] {
			case "e621":
				booru = api.E621
			case "e926":
				booru = api.E926
			}
		case "-t":
			fallthrough
		case "--tags":
			i++
			if i >= len(args) {
				log.Fatal("Missing tags after --tags")
			}
			initialTags = strings.Split(args[i], " ")
		}
	}
}

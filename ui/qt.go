package ui

import (
	q "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/qml"
)

func Spawn() {
	q.NewQApplication([]string{})
	engine := qml.NewQQmlApplicationEngine()
	url := q.QUrl_FromLocalFile("ui/main.qml")
	q.QApplication_Exec()
}


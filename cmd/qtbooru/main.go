package main

import (
	"os"

	q "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/qml"
)

func Main() {
	q.NewQApplication(os.Args)
	engine := qml.NewQQmlApplicationEngine()
	engine.Load(q.NewQUrl3("qrc:/main.qml"))
	q.QApplication_Exec()
}


package ui

import (
	"os"

	q "github.com/mappu/miqt/qt6"
)

func Spawn() {
	q.NewQApplication([]string{})
	q.QApplication_Exec()
}


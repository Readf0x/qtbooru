package main

//go:generate miqt-rcc -Input "ui/resources.qrc" -OutputGo "resources.go" -OutputRcc "resources.rcc"

import (
	"embed"

	"github.com/mappu/miqt/qt"
)

//go:embed resources.rcc
var _resourceRcc []byte

func init() {
	_ = embed.FS{}
	qt.QResource_RegisterResourceWithRccData(&_resourceRcc[0])
}

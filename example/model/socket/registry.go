package model

import (
	"log"

	"github.com/coretrix/hitrix/service/component/socket"
)

const (
	DefaultWebsocketNamespace = "default"
)

func RegisterSocketHandler(_ *socket.Socket) {
	log.Println("Register Socket")
}

func UnRegisterSocketHandler(_ *socket.Socket) {
	log.Println("UnRegister Socket")
}

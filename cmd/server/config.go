package main

import (
	"github.com/ben-rw/ragin-mages/internal/room"
	"html/template"
)

type config struct {
	RoomReg      room.RoomRegistry
	templates    *template.Template
	Port         string
	FilepathRoot string
	URLRoot      string
}

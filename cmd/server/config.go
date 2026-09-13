package main

import (
	"github.com/ben-rw/ragin-mages/internal/room"
	"html/template"
)

type config struct {
	Port         string
	FilepathRoot string
	URLRoot      string
	RoomReg      room.RoomRegistry
	templates    *template.Template
}

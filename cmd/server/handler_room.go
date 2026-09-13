package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/google/uuid"
)

// generate 4 letter code for room, set creating player as host, redirect user to room at /room/{roomID}
func (cfg *config) handlerCreateRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{CreateUsernameError: "player must provide a username"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
		return
	}

	//gen roomID + ensure there isn't an active room with the same roomID
	var roomID string
	for {
		roomID = room.GenerateRoomID()
		_, ok := cfg.RoomReg.Get(roomID)
		if !ok {
			break
		}
	}

	newRoom := &room.Room{
		ID:                roomID,
		Players:           []*room.Player{},
		PlayerSpriteIndex: room.NewPlayerSpriteIndex(),
		Scene:             protocol.LobbyScene,
		Mu:                &sync.RWMutex{},
	}

	cfg.RoomReg.Set(roomID, newRoom)

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Path:     "/",
		HttpOnly: true,
	})

	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("couldn't create uuid: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
	})

	fmt.Printf("Creating new room at %v for %v\n", "/room/"+roomID, name)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

// parses url for roomID, creates Player with Player.Host set to false, adds Player to room.Players, redirects to /room/{roomID}
func (cfg *config) handlerJoinRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinUsernameError: "player must provide a username"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
			log.Printf("unable to serve room page with error message: %v", err)
		}
		return
	}

	roomID := r.FormValue("roomID")
	roomID = strings.ToUpper(roomID)
	currentRoom, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinCodeError: "no room with that code"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
			log.Printf("unable to serve room page with error message: %v\n", err)
		}
		return
	}

	if len(currentRoom.Players) >= room.MaxPlayers {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinRoomFullError: "room is full"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
			log.Printf("unable to serve room page with error message: %v\n", err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("couldn't create uuid: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

// parses url for roomID, serves room page html at /room/{roomID}
func (cfg *config) handlerServeRoomPage(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")

	roomID = strings.ToUpper(roomID)
	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		http.Error(w, "unable to serve room page: unregistered room code", http.StatusInternalServerError)
		return
	}

	err := cfg.templates.ExecuteTemplate(w, "room.html", roomID)
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
		log.Printf("unable to serve room page: couldn't execute template: %v\n", err)
		return
	}
}

func (cfg *config) handlerServeTestRoom(w http.ResponseWriter, r *http.Request) {
	roomID := "TEST"
	username := "tester"

	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("couldn't create uuid: %v", err)
	}

	currentRoom, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		psi := room.NewPlayerSpriteIndex()
		newRoom := &room.Room{
			ID:                roomID,
			Players:           []*room.Player{},
			Scene:             protocol.LobbyScene,
			PlayerSpriteIndex: psi,
			Mu:                &sync.RWMutex{},
		}
		cfg.RoomReg.Set(roomID, newRoom)

		err := newRoom.AppendPlayer(&room.Player{
			ID:          id.String(),
			Name:        username,
			Score:       0,
			Host:        true,
			SpriteIndex: newRoom.AssignPlayerSprite(),
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		err := currentRoom.AppendPlayer(&room.Player{
			ID:          id.String(),
			Name:        username,
			Score:       0,
			Host:        false,
			SpriteIndex: currentRoom.AssignPlayerSprite(),
		})
		if err != nil {
			log.Println(err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
	})

	err = cfg.templates.ExecuteTemplate(w, "room.html", "TEST")
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
		log.Printf("unable to serve room page: couldn't execute template: %v\n", err)
		return
	}
}

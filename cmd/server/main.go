package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ben-rw/ragin-mages/frontend"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	templates, err := frontend.LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		Port:         os.Getenv("PORT"),
		FilepathRoot: os.Getenv("FILEPATH_ROOT"),
		URLRoot:      os.Getenv("URL_ROOT"),
		RoomReg: room.RoomRegistry{
			Mu:          &sync.RWMutex{},
			ActiveRooms: map[string]*room.Room{},
		},
		templates: templates,
	}

	mux := http.NewServeMux()

	mux.Handle("GET /{$}", noCacheMiddleware(http.HandlerFunc(cfg.handlerServeLandingPage)))

	staticFS, err := fs.Sub(frontend.FS, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("POST /room", cfg.handlerCreateRoom)
	mux.HandleFunc("POST /room/join", cfg.handlerJoinRoom)

	mux.HandleFunc("GET /room/{roomID}", cfg.handlerServeRoomPage)
	mux.HandleFunc("GET /room/test", cfg.handlerServeTestRoom)

	mux.HandleFunc("GET /room/{roomID}/ws", cfg.handlerWebsocket)
	mux.HandleFunc("GET /room/test/ws", cfg.handlerWebsocket)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Serving on: %s:%s/\n", cfg.URLRoot, cfg.Port)
	log.Fatal(server.ListenAndServe())
}

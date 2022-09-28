package router

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Main struct {
	R *mux.Router
}

var main *Main

func GetMain() (*Main, error) {
	if main == nil {
		return nil, errors.New("SetRouter must be called before the GetMain function")
	}
	return main, nil
}

func SetRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the chat app"))
	}).Methods("GET")

	main = &Main{
		R: r,
	}
}

func (m *Main) AddSubRouters() {
	AddAuthSubRouter()
}

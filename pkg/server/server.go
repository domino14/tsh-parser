package server

import (
	"net/http"
)

// an http server
type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("bad content type"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(`{"msg": "ok"}`))
}

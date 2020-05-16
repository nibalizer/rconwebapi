package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Config holds configuration for the server
type Config struct {
	BindAddress string
}

type rconRequest struct {
	Address  string
	Password string
	Command  string
}

type rconResponse struct {
	Output string
}

type rconReqBody struct {
	RconRequest rconRequest
}

type rconResponseBody struct {
	RconResponse rconResponse
}

// Server manages server state
type Server struct {
	config *Config
}

func logRequest(req *http.Request) {
	log.Printf("Got %q request from %q for %q", req.Method, req.RemoteAddr, req.URL)
}

func invalidMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}

// NewServer creates a new server
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

// Run starts the Server
func (s *Server) Run() {
	s.setupHandlers()
	log.Fatal(http.ListenAndServe(s.config.BindAddress, nil))
}

func (s *Server) setupHandlers() {
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/rcon", s.rconHandler)
}

func (s *Server) indexHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	if req.Method != http.MethodGet {
		invalidMethod(w)
		return
	}

	w.Write([]byte("Hello!"))
}

func (s *Server) rconHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	if req.Method != http.MethodPost {
		invalidMethod(w)
		return
	}

	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(([]byte("Ivalid Content-Type, only application/json allowed.")))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var reqBody rconReqBody
	if err := decoder.Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request, unable to parse request body."))
		return
	}

	resp := s.makeRconRequest(&reqBody.RconRequest)
	respBody := rconResponseBody{
		RconResponse: rconResponse{
			Output: resp,
		},
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(respBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error, failed to marshall response."))
	}
}

func (s *Server) makeRconRequest(rconReq *rconRequest) string {
	return ""
}

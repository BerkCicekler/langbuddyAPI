package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store Store
}

func NewAPIServer(addr string, store Store) *APIServer {
	return &APIServer{addr: addr, store: store}
}

// initialize the router
// register the services and dependices
func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// register services
	userService := NewUsersService(&s.store)
	userService.RegisterRouters(subRouter)

	userUpdaterService := NewUserUpdaterService(&s.store)
	userUpdaterService.RegisterRouters(subRouter)

	searchService := NewSearchService(&s.store)
	searchService.RegisterRouters(subRouter)

	friendsService := NewFriendsService(&s.store)
	friendsService.RegisterRouters(subRouter)

	chatService := NewChatService(&s.store)
	chatService.RegisterRouters(subRouter)

	log.Println("Starting the API server at", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}
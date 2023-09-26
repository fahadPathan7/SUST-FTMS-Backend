package router

import (
	"ftms/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter() // create mux router. it will be used to register routes

	// register routes. we are calling controller functions for each route
	// router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	// router.HandleFunc("/api/movie/{id}", controller.GetOneMovie).Methods("GET")
	router.HandleFunc("/api/dept", controller.InsertNewDept).Methods("POST") // to insert new dept
	router.HandleFunc("/api/player", controller.InsertNewPlayer).Methods("POST") // to insert new player
	router.HandleFunc("/api/team", controller.InsertNewTeam).Methods("POST") // to insert new team
	router.HandleFunc("/api/tournament", controller.InsertNewTournament).Methods("POST") // to insert new tournament
	router.HandleFunc("/api/referee", controller.InsertNewReferee).Methods("POST") // to insert new referee
	router.HandleFunc("/api/match", controller.InsertNewMatch).Methods("POST") // to insert new match
	router.HandleFunc("/api/tiebreaker", controller.InsertNewTiebreaker).Methods("POST") // to insert new tiebreaker
	// router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controller.DeleteOneMovie).Methods("DELETE")
	// router.HandleFunc("/api/dmovies", controller.DeleteAllMovies).Methods("DELETE")

	return router
}
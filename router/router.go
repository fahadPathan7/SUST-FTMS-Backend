package router

import (
	"ftms/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter() // create mux router. it will be used to register routes

	// POST operations
	router.HandleFunc("/api/dept", controller.InsertNewDept).Methods("POST") // to insert new dept
	router.HandleFunc("/api/player", controller.InsertNewPlayer).Methods("POST") // to insert new player
	router.HandleFunc("/api/team", controller.InsertNewTeam).Methods("POST") // to insert new team
	router.HandleFunc("/api/tournament", controller.InsertNewTournament).Methods("POST") // to insert new tournament
	router.HandleFunc("/api/referee", controller.InsertNewReferee).Methods("POST") // to insert new referee
	router.HandleFunc("/api/match", controller.InsertNewMatch).Methods("POST") // to insert new match
	router.HandleFunc("/api/tiebreaker", controller.InsertNewTiebreaker).Methods("POST") // to insert new tiebreaker
	router.HandleFunc("/api/individualpunishment", controller.InsertNewIndividualPunishment).Methods("POST") // to insert new individualpunishment
	router.HandleFunc("/api/individualscore", controller.InsertNewIndividualScore).Methods("POST") // to insert new individualscore

	// GET operations
	router.HandleFunc("/api/depts", controller.GetAllDepts).Methods("GET") // to get all depts
	router.HandleFunc("/api/dept/{id}", controller.GetADept).Methods("GET") // to get a specific dept
	// router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controller.DeleteOneMovie).Methods("DELETE")
	// router.HandleFunc("/api/dmovies", controller.DeleteAllMovies).Methods("DELETE")

	return router
}
package router

import (
	"ftms/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter() // create mux router. it will be used to register routes

	// operator login
	router.HandleFunc("/api/operator/login", controller.Login).Methods("POST") // to login as an operator
	router.HandleFunc("/api/token/generate/{userEmail}", controller.GenerateToken).Methods("GET") // to generate token (for operator)
	router.HandleFunc("/api/token/validate", controller.IsTokenValid).Methods("GET") // to validate token

	// POST operations
	router.HandleFunc("/api/dept", controller.InsertNewDept).Methods("POST") // to insert new dept
	router.HandleFunc("/api/player", controller.InsertNewPlayer).Methods("POST") // to insert new player
	router.HandleFunc("/api/team", controller.InsertNewTeam).Methods("POST") // to insert new team
	router.HandleFunc("/api/tournament", controller.InsertNewTournament).Methods("POST") // to insert new tournament
	router.HandleFunc("/api/referee", controller.InsertNewReferee).Methods("POST") // to insert new referee
	router.HandleFunc("/api/match", controller.InsertNewMatch).Methods("POST") // to insert new match
	router.HandleFunc("/api/match/startingeleven", controller.InsertNewStartingEleven).Methods("POST") // to insert new lineup
	router.HandleFunc("/api/tiebreaker", controller.InsertNewTiebreaker).Methods("POST") // to insert new tiebreaker
	router.HandleFunc("/api/individualpunishment", controller.InsertNewIndividualPunishment).Methods("POST") // to insert new individualpunishment
	router.HandleFunc("/api/individualscore", controller.InsertNewIndividualScore).Methods("POST") // to insert new individualscore

	// GET operations
	router.HandleFunc("/api/player/{playerRegNo}", controller.GetAPlayer).Methods("GET") // to get a specific player
	router.HandleFunc("/api/depts", controller.GetAllDepts).Methods("GET") // to get all depts
	router.HandleFunc("/api/dept/{deptCode}", controller.GetADept).Methods("GET") // to get a specific dept
	router.HandleFunc("/api/tournaments", controller.GetAllTournaments).Methods("GET") // to get all tournaments
	router.HandleFunc("/api/tournament/{tournamentId}", controller.GetATournament).Methods("GET") // to get a specific tournament
	router.HandleFunc("/api/tournament/teams/{tournamentId}", controller.GetAllTeamsOfATournament).Methods("GET") // to get all teams of a tournament
	router.HandleFunc("/api/dept/players/{deptCode}", controller.GetPlayersOfADept).Methods("GET") // to get all players of a dept
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.GetATeamOfATournament).Methods("GET") // to get a team of a tournament
	router.HandleFunc("/api/tournament/matches/{tournamentId}", controller.GetAllMatchesOfATournament).Methods("GET") // to get all matches of a tournament
	router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.GetAMatchOfATournament).Methods("GET") // to get a match of a tournament
	router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{deptCode}", controller.GetStartingElevenOfATeamOfAMatch).Methods("GET") // to get starting eleven of a team of a match
	router.HandleFunc("/api/referees", controller.GetAllReferees).Methods("GET") // to get all referees
	router.HandleFunc("/api/referee/{refereeId}", controller.GetAReferee).Methods("GET") // to get a specific referee
	router.HandleFunc("/api/tournament/tiebreakers/{tournamentId}", controller.GetAllTiebreakersOfATournament).Methods("GET") // to get all tiebreakers of a tournament
	router.HandleFunc("/api/tournament/tiebreaker/{tournamentId}/{matchId}", controller.GetATiebreakerOfATournament).Methods("GET") // to get a tiebreaker (a match) of a tournament
	router.HandleFunc("/api/tournament/individualscores/{tournamentId}", controller.GetAllIndividualScoresOfATournament).Methods("GET") // to get all individualscores (all players) of a tournament
	router.HandleFunc("/api/tournament/player/individualscores/{tournamentId}/{playerRegNo}", controller.GetAllIndividualScoresOfAPlayerInATournament).Methods("GET") // to get all individualscores of a player in a tournament
	router.HandleFunc("/api/tournament/match/team/individualscores/{tournamentId}/{matchId}/{teamDeptCode}", controller.GetAllIndividualScoresOfAMatchByATeam).Methods("GET") // to get all individual scores of a match by a team
	router.HandleFunc("/api/tournament/individualpunishments/{tournamentId}", controller.GetAllIndividualPunishmentsOfATournament).Methods("GET") // to get all individualpunishments (all players) of a tournament
	router.HandleFunc("/api/tournament/match/team/individualpunishments/{tournamentId}/{matchId}/{teamDeptCode}", controller.GetAllIndividualPunishmentsOfAMatchByATeam).Methods("GET") // to get all individual punishments of a match by a team
	router.HandleFunc("/api/tournament/player/individualpunishments/{tournamentId}/{playerRegNo}", controller.GetAllIndividualPunishmentsOfAPlayerInATournament).Methods("GET") // to get all individual punishments of a player in a tournament

	// PUT operations
	router.HandleFunc("/api/player/{playerRegNo}", controller.UpdateAPlayer).Methods("PUT") // to update a player
	router.HandleFunc("/api/dept/{deptCode}", controller.UpdateADept).Methods("PUT") // to update a dept
	router.HandleFunc("/api/tournament/{tournamentId}", controller.UpdateATournament).Methods("PUT") // to update a tournament
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.UpdateATeam).Methods("PUT") // to update a team
	router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.UpdateAMatch).Methods("PUT") // to update a match
	router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}", controller.UpdateAStartingEleven).Methods("PUT") // to update a lineup
	router.HandleFunc("/api/referee/{refereeId}", controller.UpdateAReferee).Methods("PUT") // to update a referee
	router.HandleFunc("/api/match/tiebreaker/{tournamentId}/{matchId}", controller.UpdateATiebreaker).Methods("PUT") // to update a tiebreaker of a match
	router.HandleFunc("/api/match/individualscore/{tournamentId}/{matchId}/{playerRegNo}", controller.UpdateAnIndividualScore).Methods("PUT") // to update an individual score of a match
	router.HandleFunc("/api/match/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}", controller.UpdateAnIndividualPunishment).Methods("PUT") // to update an individual punishment of a match

	// DELETE operations
	router.HandleFunc("/api/match/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}", controller.DeleteAnIndividualPunishment).Methods("DELETE") // to delete an individual punishment of a match
	router.HandleFunc("/api/match/individualscore/{tournamentId}/{matchId}/{playerRegNo}", controller.DeleteAnIndividualScore).Methods("DELETE") // to delete an individual score of a match
	router.HandleFunc("/api/match/tiebreaker/{tournamentId}/{matchId}", controller.DeleteATiebreaker).Methods("DELETE") // to delete a tiebreaker of a match
	router.HandleFunc("/api/referee/{refereeId}", controller.DeleteAReferee).Methods("DELETE") // to delete a referee
	router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.DeleteAMatch).Methods("DELETE") // to delete a match
	router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}", controller.DeleteAStartingEleven).Methods("DELETE") // to delete a lineup
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.DeleteATeam).Methods("DELETE") // to delete a team
	router.HandleFunc("/api/tournament/{tournamentId}", controller.DeleteATournament).Methods("DELETE") // to delete a tournament
	router.HandleFunc("/api/dept/{deptCode}", controller.DeleteADept).Methods("DELETE") // to delete a dept
	router.HandleFunc("/api/player/{playerRegNo}", controller.DeleteAPlayer).Methods("DELETE") // to delete a player


	return router
}
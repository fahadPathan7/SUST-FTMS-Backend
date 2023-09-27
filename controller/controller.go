package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ftms/models"
	"strconv"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// connecting to mysql database
func CreateDbConnection() {
	var err error
	db, err = sql.Open("mysql", "fahadftms:fahadftms@tcp(localhost:3306)/ftms")
	// port 3306 is the default port for mysql in xampp
	// here ftms is the database name

	if err != nil {
		fmt.Println("Error connecting databse!")
		panic(err.Error())
	}

	// defer db.Close()
	fmt.Println("Successfully connected to mysql database")
}

// verifications.

// check if player exists in database
func playerExists(playerRegNo int) bool {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if team exists in database
func teamExists(tournamentId string, deptCode int) bool {
	var team models.Team
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.Player1RegNo, &team.Player2RegNo, &team.Player3RegNo, &team.Player4RegNo, &team.Player5RegNo, &team.Player6RegNo, &team.Player7RegNo, &team.Player8RegNo, &team.Player9RegNo, &team.Player10RegNo, &team.Player11RegNo, &team.Player12RegNo, &team.Player13RegNo, &team.Player14RegNo, &team.Player15RegNo, &team.Player16RegNo, &team.Player17RegNo, &team.Player18RegNo, &team.Player19RegNo, &team.Player20RegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if tournament exists in database
func tournamentExists(tournamentId string) bool {
	var tournament models.Tournament
	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if dept exists in database
func deptExists(deptCode int) bool {
	var dept models.Dept
	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if referee exists in database
func refereeExists(refereeID int) bool {
	var referee models.Referee
	err := db.QueryRow("SELECT * FROM tblreferee WHERE refereeID = ?", refereeID).Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if match exists in database
func matchExists(tournamentId string, matchId string) bool {
	var match models.Match
	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if tiebreaker exists in database
func tiebreakerExists(tournamentId string, matchId string) bool {
	var tiebreaker models.Tiebreaker
	err := db.QueryRow("SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}





// insert dept info into database
func insertNewDept(dept models.Dept) {
	// dept.DeptCode is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbldept VALUES (?, ?, ?, ?)", dept.DeptCode, dept.DeptName, dept.DeptHeadName, dept.DeptShortName)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new dept
func InsertNewDept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	if !deptExists(dept.DeptCode) {
		insertNewDept(dept)
		json.NewEncoder(w).Encode(dept)
	} else {
		json.NewEncoder(w).Encode("Dept already exists!")
	}
}





// insert player info into database
func insertNewPlayer(player models.Player) {
	// player.PlayerRegNo is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblplayer VALUES (?, ?, ?, ?, ?)", player.PlayerRegNo, player.PlayerSession, player.PlayerSemester, player.PlayerName, player.PlayerDeptCode)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new player
func InsertNewPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	// check if player already exists
	if playerExists(player.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player already exists!")
		return
	}

	// check if dept exists
	if !deptExists(player.PlayerDeptCode) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// insert new player
	insertNewPlayer(player)
	json.NewEncoder(w).Encode(player)
}





// insert team info into database
func insertNewTeam(team models.Team) {
	// team.TournamentId is int type. and team.deptCode is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblteam VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", team.TournamentId, team.TeamSubmissionDate, team.DeptCode, team.TeamManager, team.TeamCaptainRegID, team.Player1RegNo, team.Player2RegNo, team.Player3RegNo, team.Player4RegNo, team.Player5RegNo, team.Player6RegNo, team.Player7RegNo, team.Player8RegNo, team.Player9RegNo, team.Player10RegNo, team.Player11RegNo, team.Player12RegNo, team.Player13RegNo, team.Player14RegNo, team.Player15RegNo, team.Player16RegNo, team.Player17RegNo, team.Player18RegNo, team.Player19RegNo, team.Player20RegNo)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
}

// return player's dept code from tblplayer in database
func getPlayerDeptCode(playerRegNo int) int {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		panic(err.Error())
	}

	return player.PlayerDeptCode
}

// controller function to insert new team
func InsertNewTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)

	// check if team already exists
	if teamExists(team.TournamentId, team.DeptCode) {
		json.NewEncoder(w).Encode("Team already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(team.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if dept exists
	if !deptExists(team.DeptCode) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// check if all players exist
	if !playerExists(team.TeamCaptainRegID) {
		json.NewEncoder(w).Encode("Team captain doesn't exist!")
		return
	}
	if !playerExists(team.Player1RegNo) {
		json.NewEncoder(w).Encode("Player1 doesn't exist!")
		return
	}
	if !playerExists(team.Player2RegNo) {
		json.NewEncoder(w).Encode("Player2 doesn't exist!")
		return
	}
	if !playerExists(team.Player3RegNo) {
		json.NewEncoder(w).Encode("Player3 doesn't exist!")
		return
	}
	if !playerExists(team.Player4RegNo) {
		json.NewEncoder(w).Encode("Player4 doesn't exist!")
		return
	}
	if !playerExists(team.Player5RegNo) {
		json.NewEncoder(w).Encode("Player5 doesn't exist!")
		return
	}
	if !playerExists(team.Player6RegNo) {
		json.NewEncoder(w).Encode("Player6 doesn't exist!")
		return
	}
	if !playerExists(team.Player7RegNo) {
		json.NewEncoder(w).Encode("Player7 doesn't exist!")
		return
	}
	if !playerExists(team.Player8RegNo) {
		json.NewEncoder(w).Encode("Player8 doesn't exist!")
		return
	}
	if !playerExists(team.Player9RegNo) {
		json.NewEncoder(w).Encode("Player9 doesn't exist!")
		return
	}
	if !playerExists(team.Player10RegNo) {
		json.NewEncoder(w).Encode("Player10 doesn't exist!")
		return
	}
	if !playerExists(team.Player11RegNo) {
		json.NewEncoder(w).Encode("Player11 doesn't exist!")
		return
	}
	if !playerExists(team.Player12RegNo) {
		json.NewEncoder(w).Encode("Player12 doesn't exist!")
		return
	}
	if !playerExists(team.Player13RegNo) {
		json.NewEncoder(w).Encode("Player13 doesn't exist!")
		return
	}
	if !playerExists(team.Player14RegNo) {
		json.NewEncoder(w).Encode("Player14 doesn't exist!")
		return
	}
	if !playerExists(team.Player15RegNo) {
		json.NewEncoder(w).Encode("Player15 doesn't exist!")
		return
	}
	if !playerExists(team.Player16RegNo) {
		json.NewEncoder(w).Encode("Player16 doesn't exist!")
		return
	}
	if !playerExists(team.Player17RegNo) {
		json.NewEncoder(w).Encode("Player17 doesn't exist!")
		return
	}
	if !playerExists(team.Player18RegNo) {
		json.NewEncoder(w).Encode("Player18 doesn't exist!")
		return
	}
	if !playerExists(team.Player19RegNo) {
		json.NewEncoder(w).Encode("Player19 doesn't exist!")
		return
	}
	if !playerExists(team.Player20RegNo) {
		json.NewEncoder(w).Encode("Player20 doesn't exist!")
		return
	}

	// check if all players are from same dept.
	if getPlayerDeptCode(team.TeamCaptainRegID) != team.DeptCode {
		json.NewEncoder(w).Encode("Team captain is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player1RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player1 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player2RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player2 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player3RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player3 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player4RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player4 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player5RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player5 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player6RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player6 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player7RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player7 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player8RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player8 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player9RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player9 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player10RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player10 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player11RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player11 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player12RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player12 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player13RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player13 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player14RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player14 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player15RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player15 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player16RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player16 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player17RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player17 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player18RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player18 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player19RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player19 is not from this dept!")
		return
	}
	if getPlayerDeptCode(team.Player20RegNo) != team.DeptCode {
		json.NewEncoder(w).Encode("Player20 is not from this dept!")
		return
	}


	// insert new team
	insertNewTeam(team)
	json.NewEncoder(w).Encode(team)
}





// insert tournament info into database
func insertNewTournament(tournament models.Tournament) {
	// tournament.TournamentId is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbltournament VALUES(?, ?, ?)", tournament.TournamentId, tournament.TournamentName, tournament.TournamentYear)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new tournament
func InsertNewTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var tournament models.Tournament
	_ = json.NewDecoder(r.Body).Decode(&tournament)

	if !tournamentExists(tournament.TournamentId) {
		insertNewTournament(tournament)
		json.NewEncoder(w).Encode(tournament)
	} else {
		json.NewEncoder(w).Encode("Tournament already exists!")
	}
}





// insert referee info into database
func insertNewReferee(referee models.Referee) {
	// referee.RefereeID is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblreferee VALUES (?, ?, ?)", referee.RefereeID, referee.RefereeName, referee.RefereeInstitute)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new referee
func InsertNewReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var referee models.Referee
	_ = json.NewDecoder(r.Body).Decode(&referee)

	if refereeExists(referee.RefereeID) {
		json.NewEncoder(w).Encode("Referee already exists!")
		return
	}

	insertNewReferee(referee)
	json.NewEncoder(w).Encode(referee)
}





// insert match info into database
func insertNewMatch(match models.Match) {
	// match.TournamentId is int type. and match.MatchId is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblmatch VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", match.TournamentId, match.MatchId, match.MatchDate, match.Team1DeptCode, match.Team2DeptCode, match.Team1Score, match.Team2Score, match.WinnerTeamDeptCode, match.MatchRefereeID, match.MatchLinesman1ID, match.MatchLinesman2ID, match.MatchFourthRefereeID)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new match
func InsertNewMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var match models.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// check if match already exists
	if matchExists(match.TournamentId, match.MatchId) {
		json.NewEncoder(w).Encode("Match already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(match.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if team1 exists
	if !teamExists(match.TournamentId, match.Team1DeptCode) {
		json.NewEncoder(w).Encode("Team1 doesn't exist!")
		return
	}

	// check if team2 exists
	if !teamExists(match.TournamentId, match.Team2DeptCode) {
		json.NewEncoder(w).Encode("Team2 doesn't exist!")
		return
	}

	// check if referee exists
	if !refereeExists(match.MatchRefereeID) {
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman1ID) {
		json.NewEncoder(w).Encode("Linesman1 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman2ID) {
		json.NewEncoder(w).Encode("Linesman2 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchFourthRefereeID) {
		json.NewEncoder(w).Encode("Fourth referee doesn't exist!")
		return
	}

	// insert new match
	insertNewMatch(match)
	json.NewEncoder(w).Encode(match)
}





// insert tiebreaker info into database
func insertNewTiebreaker(tiebreaker models.Tiebreaker) {
	// tiebreaker.TournamentId is int type. and tiebreaker.MatchId is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tbltiebreaker VALUES (?, ?, ?, ?, ?, ?)", tiebreaker.TournamentId, tiebreaker.MatchId, tiebreaker.Team1DeptCode, tiebreaker.Team2DeptCode, tiebreaker.Team1TieBreakerScore, tiebreaker.Team2TieBreakerScore)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new tiebreaker
func InsertNewTiebreaker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var tiebreaker models.Tiebreaker
	_ = json.NewDecoder(r.Body).Decode(&tiebreaker)

	// check if tournament exists
	if !tournamentExists(tiebreaker.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(tiebreaker.TournamentId, tiebreaker.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team1 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team1DeptCode) {
		json.NewEncoder(w).Encode("Team1 doesn't exist!")
		return
	}

	// check if team2 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team2DeptCode) {
		json.NewEncoder(w).Encode("Team2 doesn't exist!")
		return
	}

	// insert new tiebreaker
	insertNewTiebreaker(tiebreaker)
	json.NewEncoder(w).Encode(tiebreaker)
}





// insert individual score info into database
func insertNewIndividualScore(individualScore models.IndividualScore) {
	// individualScore.TournamentId is int type. and individualScore.MatchId is int type. and individualScore.PlayerRegNo is int type. and all are primary key.
	insert, err := db.Query("INSERT INTO tblindividualscore VALUES (?, ?, ?, ?, ?)", individualScore.TournamentId, individualScore.MatchId, individualScore.PlayerRegNo, individualScore.TeamDeptCode, individualScore.Goals)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new individual score
func InsertNewIndividualScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var individualScore models.IndividualScore
	_ = json.NewDecoder(r.Body).Decode(&individualScore)

	// check if tournament exists
	if !tournamentExists(individualScore.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualScore.TournamentId, individualScore.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualScore.TournamentId, individualScore.TeamDeptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if player exists
	if !playerExists(individualScore.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// insert new individual score
	insertNewIndividualScore(individualScore)
	json.NewEncoder(w).Encode(individualScore)
}





// insert individual punishment info into database
func insertNewIndividualPunishment(individualPunishment models.IndividualPunishment) {
	// individualPunishment.TournamentId is int type. and individualPunishment.MatchId is int type. and individualPunishment.PlayerRegNo is int type. and all are primary key.
	insert, err := db.Query("INSERT INTO tblindividualpunishment VALUES (?, ?, ?, ?, ?)", individualPunishment.TournamentId, individualPunishment.MatchId, individualPunishment.PlayerRegNo, individualPunishment.TeamDeptCode, individualPunishment.PunishmentType)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new individual punishment
func InsertNewIndividualPunishment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var individualPunishment models.IndividualPunishment
	_ = json.NewDecoder(r.Body).Decode(&individualPunishment)

	// check if tournament exists
	if !tournamentExists(individualPunishment.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualPunishment.TournamentId, individualPunishment.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualPunishment.TournamentId, individualPunishment.TeamDeptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if player exists
	if !playerExists(individualPunishment.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// insert new individual punishment
	insertNewIndividualPunishment(individualPunishment)
	json.NewEncoder(w).Encode(individualPunishment)
}










// getting info from database

// get all depts from database
func getAllDepts() []models.Dept {
	var dept models.Dept
	var depts []models.Dept

	result, err := db.Query("SELECT * FROM tbldept")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

		if err != nil {
			panic(err.Error())
		}

		depts = append(depts, dept)
	}

	return depts
}

// controller function to get all depts
func GetAllDepts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var depts []models.Dept
	depts = getAllDepts()

	json.NewEncoder(w).Encode(depts)
}





// get a dept from database
func getADept(deptCode int) models.Dept {
	var dept models.Dept

	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

	if err != nil {
		//panic(err.Error())
		return models.Dept{}
	}

	return dept
}

// controller function to get a dept
func GetADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/dept/{deptCode}", controller.GetADept).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	var dept models.Dept
	dept = getADept(id)

	json.NewEncoder(w).Encode(dept)
}





// get all the tournaments
func getAllTournaments() []models.Tournament {
	var tournament models.Tournament
	var tournaments []models.Tournament

	result, err := db.Query("SELECT * FROM tbltournament")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

		if err != nil {
			panic(err.Error())
		}

		tournaments = append(tournaments, tournament)
	}

	return tournaments
}

// controller function to get all tournaments
func GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	var tournaments []models.Tournament
	tournaments = getAllTournaments()

	json.NewEncoder(w).Encode(tournaments)
}





// get a tournament
func getATournament(tournamentId string) models.Tournament {
	var tournament models.Tournament

	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

	if err != nil {
		//panic(err.Error())
		return models.Tournament{}
	}

	return tournament
}

// controller function to get a tournament
func GetATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/{tournamentId}", controller.GetATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var tournament models.Tournament
	tournament = getATournament(id)

	json.NewEncoder(w).Encode(tournament)
}





// get all teams of a tournament
func getAllTeamsOfATournament(tournamentId string) []models.Team {
	var team models.Team
	var teams []models.Team

	result, err := db.Query("SELECT * FROM tblteam WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.Player1RegNo, &team.Player2RegNo, &team.Player3RegNo, &team.Player4RegNo, &team.Player5RegNo, &team.Player6RegNo, &team.Player7RegNo, &team.Player8RegNo, &team.Player9RegNo, &team.Player10RegNo, &team.Player11RegNo, &team.Player12RegNo, &team.Player13RegNo, &team.Player14RegNo, &team.Player15RegNo, &team.Player16RegNo, &team.Player17RegNo, &team.Player18RegNo, &team.Player19RegNo, &team.Player20RegNo)

		if err != nil {
			panic(err.Error())
		}

		teams = append(teams, team)
	}

	return teams
}

// controller function to get all teams of a tournament
func GetAllTeamsOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/teams/{tournamentId}", controller.GetAllTeamsOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var teams []models.Team
	teams = getAllTeamsOfATournament(id)

	json.NewEncoder(w).Encode(teams)
}





// get players of a dept
func getPlayersOfADept(deptCode int) []models.Player {
	var player models.Player
	var players []models.Player

	result, err := db.Query("SELECT * FROM tblplayer WHERE playerDeptCode = ?", deptCode)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

		if err != nil {
			panic(err.Error())
		}

		players = append(players, player)
	}

	return players
}

// controller function to get players of a dept
func GetPlayersOfADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/dept/players/{deptCode}", controller.GetPlayersOfADept).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	var players []models.Player
	players = getPlayersOfADept(id)

	json.NewEncoder(w).Encode(players)
}





// get a team of a tournament
func getATeamOfATournament(tournamentId string, deptCode int) models.Team {
	var team models.Team

	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.Player1RegNo, &team.Player2RegNo, &team.Player3RegNo, &team.Player4RegNo, &team.Player5RegNo, &team.Player6RegNo, &team.Player7RegNo, &team.Player8RegNo, &team.Player9RegNo, &team.Player10RegNo, &team.Player11RegNo, &team.Player12RegNo, &team.Player13RegNo, &team.Player14RegNo, &team.Player15RegNo, &team.Player16RegNo, &team.Player17RegNo, &team.Player18RegNo, &team.Player19RegNo, &team.Player20RegNo)

	if err != nil {
		//panic(err.Error())
		return models.Team{}
	}

	return team
}

// controller function to get a team of a tournament
func GetATeamOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.GetATeamOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	tournamentId, _ := params["tournamentId"]

	deptCode, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// team exists or not
	if !teamExists(tournamentId, deptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	var team models.Team
	team = getATeamOfATournament(tournamentId, deptCode)

	json.NewEncoder(w).Encode(team)
}





// get all matches of a tournament
func getAllMatchesOfATournament(tournamentId string) []models.Match {
	var match models.Match
	var matches []models.Match

	result, err := db.Query("SELECT * FROM tblmatch WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

		if err != nil {
			panic(err.Error())
		}

		matches = append(matches, match)
	}

	return matches
}

// controller function to get all matches of a tournament
func GetAllMatchesOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/matches/{tournamentId}", controller.GetAllMatchesOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var matches []models.Match
	matches = getAllMatchesOfATournament(id)

	json.NewEncoder(w).Encode(matches)
}





// get a match of a tournament
func getAMatchOfATournament(tournamentId string, matchId string) models.Match {
	var match models.Match

	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

	if err != nil {
		//panic(err.Error())
		return models.Match{}
	}

	return match
}

// controller function to get a match of a tournament
func GetAMatchOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.GetAMatchOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]


	// match exists or not
	if !matchExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	var match models.Match
	match = getAMatchOfATournament(tournamentId, matchId)

	json.NewEncoder(w).Encode(match)
}





// get a player
func getAPlayer(playerRegNo int) models.Player {
	var player models.Player

	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		//panic(err.Error())
		return models.Player{}
	}

	return player
}

// controller function to get a player
func GetAPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/player/{playerRegNo}", controller.GetAPlayer).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["playerRegNo"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// player exists or not
	if !playerExists(id) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	var player models.Player
	player = getAPlayer(id)

	json.NewEncoder(w).Encode(player)
}





// get all referees
func getAllReferees() []models.Referee {
	var referee models.Referee
	var referees []models.Referee

	result, err := db.Query("SELECT * FROM tblreferee")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

		if err != nil {
			panic(err.Error())
		}

		referees = append(referees, referee)
	}

	return referees
}

// controller function to get all referees
func GetAllReferees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var referees []models.Referee
	referees = getAllReferees()

	json.NewEncoder(w).Encode(referees)
}





// get a referee
func getAReferee(refereeId int) models.Referee {
	var referee models.Referee

	err := db.QueryRow("SELECT * FROM tblreferee WHERE refereeID = ?", refereeId).Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

	if err != nil {
		//panic(err.Error())
		return models.Referee{}
	}

	return referee
}

// controller function to get a referee
func GetAReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/referee/{refereeId}", controller.GetAReferee).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["refereeId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// referee exists or not
	if !refereeExists(id) {
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}

	var referee models.Referee
	referee = getAReferee(id)

	json.NewEncoder(w).Encode(referee)
}





// get all tiebreakers of a tournament
func getAllTiebreakersOfATournament(tournamentId string) []models.Tiebreaker {
	var tiebreaker models.Tiebreaker
	var tiebreakers []models.Tiebreaker

	result, err := db.Query("SELECT * FROM tbltiebreaker WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

		if err != nil {
			panic(err.Error())
		}

		tiebreakers = append(tiebreakers, tiebreaker)
	}

	return tiebreakers
}

// controller function to get all tiebreakers of a tournament
func GetAllTiebreakersOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/tiebreakers/{tournamentId}", controller.GetAllTiebreakersOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var tiebreakers []models.Tiebreaker
	tiebreakers = getAllTiebreakersOfATournament(id)

	json.NewEncoder(w).Encode(tiebreakers)
}





// get a tiebreaker of a tournament
func getATiebreakerOfATournament(tournamentId string, matchId string) models.Tiebreaker {
	var tiebreaker models.Tiebreaker

	err := db.QueryRow("SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

	if err != nil {
		//panic(err.Error())
		return models.Tiebreaker{}
	}

	return tiebreaker
}

// controller function to get a tiebreaker of a tournament
func GetATiebreakerOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/tiebreaker/{tournamentId}/{matchId}", controller.GetATiebreakerOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// tiebreaker exists or not
	if !tiebreakerExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Tiebreaker doesn't exist!")
		return
	}

	var tiebreaker models.Tiebreaker
	tiebreaker = getATiebreakerOfATournament(tournamentId, matchId)

	json.NewEncoder(w).Encode(tiebreaker)
}
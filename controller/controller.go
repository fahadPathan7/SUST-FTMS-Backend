package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ftms/models"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
func teamExists(tournamentId int, deptCode int) bool {
	var team models.Team
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.DeptHeadName, &team.TeamManager, &team.TeamCaptainRegID, &team.Player1RegNo, &team.Player2RegNo, &team.Player3RegNo, &team.Player4RegNo, &team.Player5RegNo, &team.Player6RegNo, &team.Player7RegNo, &team.Player8RegNo, &team.Player9RegNo, &team.Player10RegNo, &team.Player11RegNo, &team.Player12RegNo, &team.Player13RegNo, &team.Player14RegNo, &team.Player15RegNo, &team.Player16RegNo, &team.Player17RegNo, &team.Player18RegNo, &team.Player19RegNo, &team.Player20RegNo)

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
func tournamentExists(tournamentId int) bool {
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





// insert dept info into database
func insertNewDept(dept models.Dept) {
	// dept.DeptCode is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbldept VALUES (?, ?, ?)", dept.DeptCode, dept.DeptName, dept.DeptShortName)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// check if dept exists in database
func deptExists(deptCode int) bool {
	var dept models.Dept
	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptShortName)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
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

	if !playerExists(player.PlayerRegNo) {
		insertNewPlayer(player)
		json.NewEncoder(w).Encode(player)
	} else {
		json.NewEncoder(w).Encode("Player already exists!")
	}
}





// insert team info into database
func insertNewTeam(team models.Team) {
	// team.TournamentId is int type. and team.deptCode is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblteam VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", team.TournamentId, team.TeamSubmissionDate, team.DeptCode, team.DeptHeadName, team.TeamManager, team.TeamCaptainRegID, team.Player1RegNo, team.Player2RegNo, team.Player3RegNo, team.Player4RegNo, team.Player5RegNo, team.Player6RegNo, team.Player7RegNo, team.Player8RegNo, team.Player9RegNo, team.Player10RegNo, team.Player11RegNo, team.Player12RegNo, team.Player13RegNo, team.Player14RegNo, team.Player15RegNo, team.Player16RegNo, team.Player17RegNo, team.Player18RegNo, team.Player19RegNo, team.Player20RegNo)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
}

// controller function to insert new team
func InsertNewTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)

	if !teamExists(team.TournamentId, team.DeptCode) {
		insertNewTeam(team)
		json.NewEncoder(w).Encode(team)
	} else {
		json.NewEncoder(w).Encode("Team already exists!")
	}
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
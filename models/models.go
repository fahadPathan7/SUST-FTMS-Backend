package models

type Dept struct {
	DeptCode      int    `json:"deptCode"`
	DeptName      string `json:"deptName"`
	DeptShortName string `json:"deptShortName"`
}

type Player struct {
	PlayerRegNo    int    `json:"playerRegNo"`
	PlayerName     string `json:"playerName"`
	PlayerDeptCode int    `json:"playerDeptCode"`
	PlayerSession  string `json:"playerSession"`
	PlayerSemester int    `json:"playerSemester"`
}

type Team struct {
	TournamentId      int    `json:"tournamentId"`
	TeamSumissionDate string `json:"teamSumissionDate"`
	DeptCode          int    `json:"deptCode"`
	DeptHeadName      string `json:"deptHeadName"`
	TeamManager       string `json:"teamManager"`
	TeamCaptainRegID  int    `json:"teamCaptainRegID"`
	Player1RegNo      int    `json:"player1RegNo"`
	Player2RegNo      int    `json:"player2RegNo"`
	Player3RegNo      int    `json:"player3RegNo"`
	Player4RegNo      int    `json:"player4RegNo"`
	Player5RegNo      int    `json:"player5RegNo"`
	Player6RegNo      int    `json:"player6RegNo"`
	Player7RegNo      int    `json:"player7RegNo"`
	Player8RegNo      int    `json:"player8RegNo"`
	Player9RegNo      int    `json:"player9RegNo"`
	Player10RegNo     int    `json:"player10RegNo"`
	Player11RegNo     int    `json:"player11RegNo"`
	Player12RegNo     int    `json:"player12RegNo"`
	Player13RegNo     int    `json:"player13RegNo"`
	Player14RegNo     int    `json:"player14RegNo"`
	Player15RegNo     int    `json:"player15RegNo"`
	Player16RegNo     int    `json:"player16RegNo"`
	Player17RegNo     int    `json:"player17RegNo"`
	Player18RegNo     int    `json:"player18RegNo"`
	Player19RegNo     int    `json:"player19RegNo"`
	Player20RegNo     int    `json:"player20RegNo"`
}

type IndividualPunishment struct {
	TournamentId   int    `json:"tournamentId"`
	MatchId        int    `json:"matchId"`
	PlayerRegNo    int    `json:"playerRegNo"`
	TeamDeptCode   int    `json:"teamDeptCode"`
	PunishmentType string `json:"punishmentType"`
}

type IndividualScore struct {
	TournamentId int `json:"tournamentId"`
	MatchId      int `json:"matchId"`
	PlayerRegNo  int `json:"playerRegNo"`
	TeamDeptCode int `json:"teamDeptCode"`
	Goals        int `json:"goals"`
}

type Tournament struct {
	TournamentId   int    `json:"tournamentId"`
	TournamentName string `json:"tournamentName"`
	TournamentYear string `json:"tournamentYear"`
}

type TieBreaker struct {
	TournamentId         int `json:"tournamentId"`
	MatchId              int `json:"matchId"`
	Team1DeptCode        int `json:"team1DeptCode"`
	Team2DeptCode        int `json:"team2DeptCode"`
	Team1TieBreakerScore int `json:"team1TieBreakerScore"`
	Team2TieBreakerScore int `json:"team2TieBreakerScore"`
}

type Match struct {
	TournamentId int    `json:"tournamentId"`
	MatchId      int    `json:"matchId"`
	MatchDate    string `json:"matchDate"`
	Team1DeptCode int `json:"team1DeptCode"`
	Team2DeptCode int `json:"team2DeptCode"`
	Team1Score   int    `json:"team1Score"`
	Team2Score   int    `json:"team2Score"`
	WinnerTeamDeptCode int `json:"winnerTeamDeptCode"`
	MatchRefereeID int `json:"matchRefereeID"`
	MatchLinesman1ID int `json:"matchLinesman1ID"`
	MatchLinesman2ID int `json:"matchLinesman2ID"`
	MatchFourthRefereeID int `json:"matchFourthRefereeID"`
}

type Referee struct {
	RefereeID int `json:"refereeID"`
	RefereeName string `json:"refereeName"`
	RefereeInstitute string `json:"refereeInstitute"`
}
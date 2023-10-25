# My RESTful API

This comprehensive RESTful API facilitates the management of tournaments, players, teams, referees, and other essential components.

**Note:** For PUT and POST operations, ensure that all required information is provided.


## Endpoints

### Operator

- `POST /api/operator/login` - Operator login
- `GET /api/token/generate/{userEmail}` - Generate token
- `GET /api/token/validate` - Validate token

**A JSON sample for operator**
```json
{
  "email": "YourEmail",
  "password": "YourPassword"
}
```

### Department

- `POST /api/dept` - Insert a new department
- `PUT /api/dept/{deptCode}` - Update a department
- `GET /api/depts` - Get all departments
- `GET /api/dept/{deptCode}` - Get a specific department

**A JSON sample for department**
```json
{
  "deptCode": 1,
  "deptName": "Department Name",
  "deptHeadName": "Department Head Name",
  "deptShortName": "DeptShortName"
}
```

### Player

- `POST /api/player` - Insert a new player
- `PUT /api/player/{playerRegNo}` - Update a player
- `GET /api/player/{playerRegNo}` - Get a specific player

**A JSON sample for player**
```json
{
  "playerRegNo": 1,
  "playerSession": "Player Session",
  "playerSemester": 1,
  "playerName": "Player Name",
  "playerDeptCode": 1,
  "playerJerseyNo": 10
}
```

### Team

- `POST /api/team` - Insert a new team
- `PUT /api/tournament/team/{tournamentId}/{deptCode}` - Update a team
- `GET /api/tournament/team/{tournamentId}/{deptCode}` - Get a team in a tournament

**A JSON sample for team**
```json
{
  "tournamentId": "Tournament ID",
  "teamSubmissionDate": "Submission Date",
  "deptCode": 1,
  "teamManager": "Team Manager",
  "teamCaptainRegID": 1,
  "playerRegNo": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
                  11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
  "isKnockedOut": false
}
```

### Tournament

- `POST /api/tournament` - Insert a new tournament
- `PUT /api/tournament/{tournamentId}` - Update a tournament
- `GET /api/tournaments` - Get all tournaments
- `GET /api/tournament/{tournamentId}` - Get a specific tournament

**A JSON sample for tournament**
```json
{
  "tournamentId": "Tournament ID",
  "tournamentName": "Tournament Name",
  "startingDate": "Starting Date",
  "endingDate": "Ending Date"
}
```

### Referee

- `POST /api/referee` - Insert a new referee
- `PUT /api/referee/{refereeId}` - Update a referee
- `GET /api/referees` - Get all referees
- `GET /api/referee/{refereeId}` - Get a specific referee

**A JSON sample for referee**
```json
{
  "refereeID": 1,
  "refereeName": "Referee Name",
  "refereeInstitute": "Referee Institute"
}
```

### Match

- `POST /api/match` - Insert a new match
- `PUT /api/match/{tournamentId}/{matchId}` - Update a match
- `DELETE /api/match/{tournamentId}/{matchId}` - Delete a match

**A JSON sample for match**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "matchDate": "Match Date",
  "team1DeptCode": 1,
  "team2DeptCode": 2,
  "team1Score": 3,
  "team2Score": 4,
  "winnerTeamDeptCode": 1,
  "matchRefereeID": 5,
  "matchLinesman1ID": 6,
  "matchLinesman2ID": 7,
  "matchFourthRefereeID": 8
}
```

### Starting Eleven

- `POST /api/match/startingeleven` - Insert a new starting eleven
- `PUT /api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}` - Update a starting eleven
- `DELETE /api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}` - Delete a starting eleven

**A JSON sample for starting eleven**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "teamDeptCode": 1,
  "startingPlayerRegNo": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11],
  "substitutePlayerRegNo": [12, 13, 14],
  "substituedPlayerRegNo": [1, 2, 3]
}
```

### Individual Punishment

- `POST /api/individualpunishment` - Insert a new individual punishment
- `PUT /api/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}` - Update an individual punishment
- `DELETE /api/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}` - Delete an individual punishment

**A JSON sample for individual punishment**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "playerRegNo": 1,
  "teamDeptCode": 1,
  "punishmentType": "Punishment Type"
}
```

### Individual Score

- `POST /api/individualscore` - Insert a new individual score
- `PUT /api/individualscore/{tournamentId}/{matchId}/{playerRegNo}` - Update an individual score
- `DELETE /api/individualscore/{tournamentId}/{matchId}/{playerRegNo}` - Delete an individual score

**A JSON sample for individual score**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "playerRegNo": 1,
  "teamDeptCode": 1,
  "goals": 2
}
```

### Tiebreaker

- `POST /api/tiebreaker` - Insert a new tiebreaker
- `PUT /api/tiebreaker/{tournamentId}/{matchId}` - Update a tiebreaker
- `DELETE /api/tiebreaker/{tournamentId}/{matchId}` - Delete a tiebreaker

**A JSON sample for tiebreaker**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "team1DeptCode": 1,
  "team2DeptCode": 2,
  "team1TieBreakerScore": 3,
  "team2TieBreakerScore": 4
}
```

## Setting Up the Backend

**1. Database Setup**

- Create a database named `ftms` in MySQL (XAAMP).

**2. Table Creation**

- Copy the codes inside `allTables.txt` file from the `Database` folder.
- Paste the contents of `allTables.sql` into your MySQL client (XAAMP) and execute the script to create the necessary tables.

**3. Running the Backend**

- Open the project in VS Code.
- Run the `main.go` file to start the backend server.

**Congratulations! You're now ready to work with the backend.**

---

Thank you for reviewing the documentation for My RESTful API. This API is designed to provide a comprehensive and user-friendly interface for managing tournaments, players, teams, referees, and other essential components for the SUST Football Tournament Management System.

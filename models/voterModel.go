package models

// registration and login of a voter
type Voter struct {
	ID             int    `json:"id"`
	Name           string `json:"name" validate:"required,max=100,min=2"`
	Age            int    `json:"age" validate:"required"`
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"required,min=6"`
	ConstituencyID int    `json:"constituency_id"`
	CandidateID    int    `json:"candidate_id"`
}

// registration of a constituency
type Constituency struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	WinnerID   int    `json:"winner_id"`
	TotalVotes int    `json:"total_votes"`
}

// candidate
type Candidate struct {
	ID             int    `json:"id"`
	Name           string `json:"name" validate:"required,max=100,min=2"`
	Age            int    `json:"age"`
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"required,min=6"`
	Constituencies []int  `json:"constituencies"`
}

// vote struct
type Vote struct {
	ID             int `json:"id"`
	VoterID        int `json:"voter_id"`
	ConstituencyID int `json:"constituency_id"`
}

// election officers login
type ElectionOfficer struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,max=20"`
}

// for logging in normal voters and candidates
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// mainly for election officers
type LoginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// for updating the constituency details
type UpdateConstituencyCredentials struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

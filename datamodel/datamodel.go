package datamodel

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Gender     string `json:"gender,omitempty"`
	Age        int    `json:"age,omitempty"`
	Profession string `json:"profession,omitempty"`
}

type Session struct {
	SessionId string `json:"session_id"`
	Email     string `json:"email"`
}

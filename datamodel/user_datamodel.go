package datamodel

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	Age        int    `json:"age"`
	Profession string `json:"profession"`
}

package models

type User struct {
	ID         int
	Username   string
	Email      string
	Password   string
	Role       string
	CreatedAt  string
	ModifiedAt string
	Active     int
	UUID       string
}

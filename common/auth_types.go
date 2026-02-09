package common

type User struct {
	Id    int
	Name  string
	Email string
}

type UserSession struct {
	Id     string
	UserId int
}

package typeauth

type User struct {
	Id          int
	Login       string
	DisplayName *string
	Email       string
}

type UserSession struct {
	Id     string
	UserId int
}

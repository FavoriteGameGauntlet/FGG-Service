package typeusers

type User struct {
	Login       string
	DisplayName *string
}

type Users = []User

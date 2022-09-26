package user

type User struct {
	FirstName string
	LastName  string
	Address   string
}

func GetUser() *User {
	return &User{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "71 Cherry Court SOUTHAMPTON SO53 5PD UK",
	}
}

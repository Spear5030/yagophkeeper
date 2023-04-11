package domain

type LoginPassword struct {
	Key      int64
	Login    string
	Password string
	Meta     string
}

type User struct {
	Email    string
	Password string
}

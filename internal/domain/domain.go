package domain

type LoginPassword struct {
	Key      int
	Login    string
	Password string
	Meta     string
}

type TextData struct {
	Key  int
	Text string
	Meta string
}

type BinaryData struct {
	Key        int
	BinaryData []byte
	Meta       string
}

type CardData struct {
	Key        int
	Number     string
	CardHolder string
	CVC        string
	Meta       string
}

type User struct {
	Email    string
	Password string
}

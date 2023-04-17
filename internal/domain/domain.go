package domain

type LoginPassword struct {
	Key      int64
	Login    string
	Password string
	Meta     string
}

type TextData struct {
	Key  int64
	Text string
	Meta string
}

type BinaryData struct {
	Key        int64
	BinaryData []byte
	Meta       string
}

type CardData struct {
	Number     string
	CardHolder string
	CVC        string
	Meta       string
}

type User struct {
	Email    string
	Password string
}

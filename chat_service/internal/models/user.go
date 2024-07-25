package models

type User struct {
	Id        int64
	FirstName string
	LastName  string
	Login     string
	Avatar    string
	PassHash  []byte
}

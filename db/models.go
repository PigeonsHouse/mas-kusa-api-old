package db

type User struct {
	Instance string `json:"instance"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}

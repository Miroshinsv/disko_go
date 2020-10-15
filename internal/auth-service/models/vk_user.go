package models

type VKUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       int    `json:"sex"`
	Nick      string `json:"screen_name"`
}

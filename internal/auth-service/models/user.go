package models

type User struct {
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Role      *int    `json:"role"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	AvatarUrl *string `json:"avatar_url"`
}

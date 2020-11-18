package models

type User struct {
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	AvatarUrl *string `json:"avatar_url"`
	VkId      *string `json:"vk_id"`
	Role      *int    `json:"role"`
}

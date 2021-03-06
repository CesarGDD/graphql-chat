// Code generated by sqlc. DO NOT EDIT.

package pg

import "time"
type AuthResponse struct {
	AuthToken *AuthToken `json:"authToken"`
	User      *User      `json:"user"`
}

type AuthToken struct {
	AccessToken string    `json:"accessToken"`
	ExpireAt    time.Time `json:"expireAt"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	ID      int `json:"id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
	User    *User  `json:"user"`
}

type NewMessageInput struct {
	UserID  int `json:"user_id"`
	Content string `json:"content"`
}

type RegisterInput struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type User struct {
	ID       int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

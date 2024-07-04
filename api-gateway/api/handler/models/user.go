package models

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Id    string `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListUsers struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

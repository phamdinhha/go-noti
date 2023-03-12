package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (u *User) GetUserId() int {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role
}

func NewFakeUser() *User {
	return &User{
		ID:       1,
		Email:    "phamdinhha95@gmail.com",
		Password: "Hapham@1",
		Name:     "Ha Pham Dinh",
		Role:     "user",
	}
}

func (u *User) FakeUserID() {
	u.ID = 1
}

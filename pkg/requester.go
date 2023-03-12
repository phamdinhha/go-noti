package pkg

const CurrentUserKey = "user"

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

package user

// User identified by unique GitHub name
type User struct {
	name      string
	resources int
}

// List Slice of users
type List []*User

// NewUser returns a user with a specified name and no resources
func NewUser(name string) *User {
	return &User{name: name}
}

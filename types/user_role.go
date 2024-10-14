package types

type UserRole string

const (
	User       UserRole = "user"
	Admin      UserRole = "admin"
	SuperAdmin UserRole = "super-admin"
)

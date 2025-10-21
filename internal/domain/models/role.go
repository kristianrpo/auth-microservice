package domain

import "fmt"

// Role represents a user role in the system
type Role string

const (
	// RoleUser is the default role for regular users (citizens)
	RoleUser Role = "USER"

	// RoleAdmin is the role for administrators
	RoleAdmin Role = "ADMIN"
)

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

// ParseRole parses a string into a Role
func ParseRole(s string) (Role, error) {
	role := Role(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return role, nil
}

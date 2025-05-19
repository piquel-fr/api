package auth

import (
	repository "github.com/PiquelChips/piquel.fr/database/generated"
)

type PolicyConfiguration struct {
	Permissions map[string]*Permission
	Roles       Roles
}

type Permission struct {
	Action     string
	Conditions Conditions
	Preset     string
}

type Conditions []func(request *Request) error

type Roles map[string]*struct {
	Name        string
	Color       string
	Permissions map[string][]*Permission
	Parents     []string
}

type Request struct {
	User      *repository.User
	Ressource Resource
	Actions   []string
}

type Resource interface {
	GetResourceName() string
	GetOwner() string
}

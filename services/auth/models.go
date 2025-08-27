package auth

import (
	"context"

	"github.com/piquel-fr/api/database/repository"
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
	Context   context.Context
}

type Resource interface {
	GetResourceName() string
	GetOwner() int32
}

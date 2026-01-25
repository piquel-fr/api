package auth

import (
	"context"

	"github.com/piquel-fr/api/database/repository"
)

type PolicyConfiguration struct {
	Presets map[string]*Permission `json:"presets"`
	Roles   Roles                  `json:"roles"`
}

type Permission struct {
	Action     string     `json:"action"`
	Conditions Conditions `json:"-"`
	Preset     string     `json:"preset"`
}

type Conditions []func(request *Request) error

type Roles map[string]*struct {
	Name        string                   `json:"name"`
	Color       string                   `json:"color"`
	Permissions map[string][]*Permission `json:"permissions"`
	Parents     []string                 `json:"parents"`
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

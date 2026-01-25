package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/utils/errors"
)

type PolicyConfiguration struct {
	Presets map[string]*Permission `json:"presets"`
	Roles   map[string]*Role       `json:"roles"`
}

func (p *PolicyConfiguration) ValidateRole(role string) error {
	_, ok := p.Roles[role]
	if !ok {
		return errors.NewError(fmt.Sprintf("role %s does not exist in current policy", role), http.StatusBadRequest)
	}
	return nil
}

type Permission struct {
	Action     string     `json:"action"`
	Conditions Conditions `json:"-"`
	Preset     string     `json:"preset"`
}

type Conditions []func(request *AuthRequest) error

type Role struct {
	Name        string                   `json:"name"`
	Color       string                   `json:"color"`
	Permissions map[string][]*Permission `json:"permissions"`
	Parents     []string                 `json:"parents"`
}

type AuthRequest struct {
	User      *repository.User
	Ressource Resource
	Actions   []string
	Context   context.Context
}

type Resource interface {
	GetResourceName() string
	GetOwner() int32
}

package handlers

import (
	"context"

	"github.com/piquel-fr/api/api"
)

type ApiHandler struct{}

func (h *ApiHandler) GetUsers(ctx context.Context, request api.GetUsersRequestObject) (api.GetUsersResponseObject, error) {
	return nil, nil
}

func (h *ApiHandler) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	return nil, nil
}

func (h *ApiHandler) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	return nil, nil
}

func (h *ApiHandler) GetUserById(ctx context.Context, request api.GetUserByIdRequestObject) (api.GetUserByIdResponseObject, error) {
	return nil, nil
}

func (h *ApiHandler) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	return nil, nil
}

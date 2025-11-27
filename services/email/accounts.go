package email

import (
	"context"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
)

type MailAccount struct {
	Unread int `json:"unread"`
	repository.MailAccount
}

type AccountInfo struct {
	MailAccount
}

func (r *realEmailService) GetAccountByEmail(ctx context.Context, email string) (repository.MailAccount, error) {
	return database.Queries.GetMailAccountByEmail(ctx, email)
}

func (r *realEmailService) ListAccounts(ctx context.Context, userId int32) ([]MailAccount, error) {
	return nil, nil
}

func (r *realEmailService) AddAccount(ctx context.Context, params repository.AddEmailAccountParams) error {
	return nil
}

func (r *realEmailService) RemoveAccount(ctx context.Context, accountId int32) error {
	return nil
}

func (r *realEmailService) GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error) {
	return AccountInfo{}, nil
}

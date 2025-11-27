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
	dbAccounts, err := database.Queries.ListUserMailAccounts(ctx, userId)
	if err != nil {
		return nil, err
	}

	accounts := []MailAccount{}
	for _, dbAccount := range dbAccounts {
		// TODO: get number of unreal emails

		account := MailAccount{
			MailAccount: dbAccount,
			Unread:      0,
		}

		// make sure we don't return username and password
		account.Username = ""
		account.Password = ""
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *realEmailService) AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error) {
	return database.Queries.AddEmailAccount(ctx, params)
}

func (r *realEmailService) RemoveAccount(ctx context.Context, accountId int32) error {
	return database.Queries.RemoveMailAccount(ctx, accountId)
}

func (r *realEmailService) GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error) {
	return AccountInfo{}, nil
}

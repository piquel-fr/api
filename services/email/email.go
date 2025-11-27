package email

import (
	"context"

	"github.com/piquel-fr/api/database/repository"
)

type EmailService interface {
	// account stuff
	GetAccountByEmail(ctx context.Context, email string) (repository.MailAccount, error)
	ListAccounts(ctx context.Context, userId int32) ([]MailAccount, error)
	AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error)
	RemoveAccount(ctx context.Context, accountId int32) error
	GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error)
}

type realEmailService struct{}

func NewRealEmailService() *realEmailService {
	return &realEmailService{}
}

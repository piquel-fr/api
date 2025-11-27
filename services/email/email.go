package email

import (
	"context"
	"fmt"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
)

type EmailService interface {
	// account stuff
	GetAccountByEmail(ctx context.Context, email string) (repository.MailAccount, error)
	ListAccounts(ctx context.Context, userId int32) ([]repository.MailAccount, error)
	AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error)
	RemoveAccount(ctx context.Context, accountId int32) error
	GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error)
}

type realEmailService struct {
	imapAddr string
}

func NewRealEmailService() *realEmailService {
	addr := fmt.Sprintf("%s:%s", config.Envs.ImapHost, config.Envs.ImapPort)

	return &realEmailService{
		imapAddr: addr,
	}
}

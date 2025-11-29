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
	CountAccounts(ctx context.Context, userId int32) (int64, error)
	AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error)
	RemoveAccount(ctx context.Context, accountId int32) error
	GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error)

	// sharing
	AddShare(ctx context.Context, params repository.AddShareParams) error
	RemoveShare(ctx context.Context, params repository.RemoveShareParams) error
	GetAccountShares(ctx context.Context, account int32) ([]int32, error)
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

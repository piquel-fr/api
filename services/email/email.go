package email

import (
	"github.com/piquel-fr/api/database/repository"
)

type EmailService interface {
	// account stuff
	GetAccountByEmail(email string) (repository.MailAccount, error)
	ListAccounts(userId int32) ([]MailAccount, error)
	AddAccount(params repository.AddEmailAccountParams) error
	RemoveAccount(accountId int32) error
	GetAccountInfo(account *repository.MailAccount) (AccountInfo, error)
}

type realEmailService struct{}

func NewRealEmailService() *realEmailService {
	return &realEmailService{}
}

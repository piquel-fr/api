package email

import "github.com/piquel-fr/api/database/repository"

type AccountInfo struct {
}

type MailAccount struct {
	repository.MailAccount
	Unread int `json:"unread"`
}

func (r *realEmailService) GetAccountByEmail(email string) (repository.MailAccount, error) {
	return repository.MailAccount{}, nil
}

func (r *realEmailService) ListAccounts(userId int32) ([]MailAccount, error) {
	return nil, nil
}

func (r *realEmailService) AddAccount(params repository.AddEmailAccountParams) error {
	return nil
}

func (r *realEmailService) RemoveAccount(accountId int32) error {
	return nil
}

func (r *realEmailService) GetAccountInfo(account *repository.MailAccount) (AccountInfo, error) {
	return AccountInfo{}, nil
}

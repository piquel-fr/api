package email

import (
	"context"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/piquel-fr/api/database/repository"
)

type AccountInfo struct {
	*repository.MailAccount
	Folders []Folder `json:"mailboxes"`
	Shares  []string `json:"shares"`
}

func (s *realEmailService) GetAccountByEmail(ctx context.Context, email string) (*repository.MailAccount, error) {
	return s.storageService.GetMailAccountByEmail(ctx, email)
}

func (s *realEmailService) ListAccounts(ctx context.Context, userId int32) ([]*repository.MailAccount, error) {
	return s.storageService.ListUserMailAccounts(ctx, userId)
}

func (s *realEmailService) CountAccounts(ctx context.Context, userId int32) (int64, error) {
	return s.storageService.CountUserMailAccounts(ctx, userId)
}

func (s *realEmailService) AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error) {
	return s.storageService.AddEmailAccount(ctx, params)
}

func (s *realEmailService) RemoveAccount(ctx context.Context, accountId int32) error {
	// TODO: remove the shares as well
	return s.storageService.DeleteMailAccount(ctx, accountId)
}

func (s *realEmailService) GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error) {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return AccountInfo{}, err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return AccountInfo{}, nil
	}

	accountInfo := AccountInfo{
		MailAccount: account,
	}

	// don't want to send sensitive data to user, for internal use only
	account.Username = ""
	account.Password = ""

	accountInfo.Folders, err = s.ListFolders(account)
	if err != nil {
		return AccountInfo{}, err
	}

	// get shares
	shares, err := s.GetAccountShares(ctx, account.ID)
	if err != nil {
		return AccountInfo{}, err
	}

	for _, share := range shares {
		user, err := s.storageService.GetUserById(ctx, share)
		if err != nil {
			return AccountInfo{}, err
		}
		accountInfo.Shares = append(accountInfo.Shares, user.Username)
	}

	return accountInfo, nil
}

func (s *realEmailService) AddShare(ctx context.Context, params repository.AddShareParams) error {
	return s.storageService.AddShare(ctx, params)
}

func (s *realEmailService) RemoveShare(ctx context.Context, userId, accountId int32) error {
	return s.storageService.DeleteShare(ctx, userId, accountId)
}

func (s *realEmailService) GetAccountShares(ctx context.Context, account int32) ([]int32, error) {
	return s.storageService.ListAccountShares(ctx, account)
}

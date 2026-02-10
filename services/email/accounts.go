package email

import (
	"context"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
)

type AccountInfo struct {
	*repository.MailAccount
	Folders []Folder `json:"mailboxes"`
	Shares  []string `json:"shares"`
}

func (r *realEmailService) GetAccountByEmail(ctx context.Context, email string) (*repository.MailAccount, error) {
	return database.Queries.GetMailAccountByEmail(ctx, email)
}

func (r *realEmailService) ListAccounts(ctx context.Context, userId int32) ([]*repository.MailAccount, error) {
	return database.Queries.ListUserMailAccounts(ctx, userId)
}

func (r *realEmailService) CountAccounts(ctx context.Context, userId int32) (int64, error) {
	return database.Queries.CountUserMailAccounts(ctx, userId)
}

func (r *realEmailService) AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error) {
	return database.Queries.AddEmailAccount(ctx, params)
}

func (r *realEmailService) RemoveAccount(ctx context.Context, accountId int32) error {
	// TODO: remove the shares as well
	return database.Queries.DeleteMailAccount(ctx, accountId)
}

func (r *realEmailService) GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error) {
	client, err := imapclient.DialTLS(r.imapAddr, nil)
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

	accountInfo.Folders, err = r.ListFolders(account)
	if err != nil {
		return AccountInfo{}, err
	}

	shares, err := r.GetAccountShares(ctx, account.ID)
	if err != nil {
		return AccountInfo{}, err
	}

	for _, share := range shares {
		user, err := database.Queries.GetUserById(ctx, share)
		if err != nil {
			return AccountInfo{}, err
		}
		accountInfo.Shares = append(accountInfo.Shares, user.Username)
	}

	return accountInfo, nil
}

func (r *realEmailService) AddShare(ctx context.Context, params repository.AddShareParams) error {
	return database.Queries.AddShare(ctx, params)
}

func (r *realEmailService) RemoveShare(ctx context.Context, userId, accountId int32) error {
	return database.Queries.DeleteShare(ctx, userId, accountId)
}

func (r *realEmailService) GetAccountShares(ctx context.Context, account int32) ([]int32, error) {
	return database.Queries.ListAccountShares(ctx, account)
}

package email

import (
	"context"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
)

type Mailbox struct {
	Name        string
	NumMessages int `json:"num_messages"`
	NumUnread   int `json:"num_unread"`
}

type AccountInfo struct {
	*repository.MailAccount
	Mailboxes []Mailbox
}

func (r *realEmailService) GetAccountByEmail(ctx context.Context, email string) (repository.MailAccount, error) {
	return database.Queries.GetMailAccountByEmail(ctx, email)
}

func (r *realEmailService) ListAccounts(ctx context.Context, userId int32) ([]repository.MailAccount, error) {
	return database.Queries.ListUserMailAccounts(ctx, userId)
}

func (r *realEmailService) AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error) {
	return database.Queries.AddEmailAccount(ctx, params)
}

func (r *realEmailService) RemoveAccount(ctx context.Context, accountId int32) error {
	return database.Queries.RemoveMailAccount(ctx, accountId)
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

	// get mailboxes
	listCmd := client.List("", "*", nil)
	defer listCmd.Close()

	for mailbox := listCmd.Next(); mailbox != nil; mailbox = listCmd.Next() {
		accountInfo.Mailboxes = append(accountInfo.Mailboxes, Mailbox{
			Name:        mailbox.Mailbox,
			NumMessages: int(*mailbox.Status.NumMessages),
			NumUnread:   int(*mailbox.Status.NumUnseen),
		})
	}

	return accountInfo, nil
}

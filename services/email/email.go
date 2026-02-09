package email

import (
	"context"
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/wneessen/go-mail"
)

type Email struct {
	To, Cc, Bcc,
	From, Sender []string
	Date    time.Time
	Subject string
}

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
	RemoveShare(ctx context.Context, params repository.DeleteShareParams) error
	GetAccountShares(ctx context.Context, account int32) ([]int32, error)

	// email stuff
	SendEmail(destination []string, from *repository.MailAccount, subject, content string) error
	CountEmailsForAccount(account *repository.MailAccount) (int, error)
	GetEmailsForAccount(account *repository.MailAccount, offset, limit int) ([]*Email, error)
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

func (r *realEmailService) SendEmail(destination []string, from *repository.MailAccount, subject, content string) error {
	message := mail.NewMsg()
	if err := message.From(from.Email); err != nil {
		return fmt.Errorf("failed to add FROM address %s: %w", from.Email, err)
	}

	for _, to := range destination {
		if err := message.AddTo(to); err != nil {
			return fmt.Errorf("failed to add TO address %s: %w", to, err)
		}
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, content)

	// Deliver the mails via SMTP
	client, err := mail.NewClient(config.Envs.SmtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(from.Username), mail.WithPassword(from.Password),
	)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %s", err)
	}
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %s", err)
	}

	return nil
}

func (r *realEmailService) CountEmailsForAccount(account *repository.MailAccount) (int, error) {
	return 0, nil
}

func (r *realEmailService) GetEmailsForAccount(account *repository.MailAccount, offset, limit int) ([]*Email, error) {
	addr := fmt.Sprintf("%s:%s", config.Envs.ImapHost, config.Envs.ImapPort)
	client, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return nil, err
	}

	// Select INBOX
	mailbox, err := client.Select("INBOX", nil).Wait()
	if err != nil {
		return nil, err
	}

	seqSet := imap.SeqSet{{Start: 1, Stop: mailbox.NumMessages}}
	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierHeader},
		},
	}

	messages, err := client.Fetch(seqSet, fetchOptions).Collect()
	if err != nil {
		return nil, err
	}

	emails := []*Email{}
	for _, msg := range messages {
		email := Email{}

		for _, to := range msg.Envelope.To {
			if to.IsGroupEnd() || to.IsGroupStart() {
				continue
			}
			email.To = append(email.To, to.Addr())
		}

		for _, cc := range msg.Envelope.Cc {
			if cc.IsGroupEnd() || cc.IsGroupStart() {
				continue
			}
			email.Cc = append(email.Cc, cc.Addr())
		}

		for _, bcc := range msg.Envelope.Bcc {
			if bcc.IsGroupEnd() || bcc.IsGroupStart() {
				continue
			}
			email.Bcc = append(email.Bcc, bcc.Addr())
		}

		for _, from := range msg.Envelope.From {
			if from.IsGroupEnd() || from.IsGroupStart() {
				continue
			}
			email.From = append(email.From, from.Addr())
		}

		for _, sender := range msg.Envelope.Sender {
			if sender.IsGroupEnd() || sender.IsGroupStart() {
				continue
			}
			email.Sender = append(email.Sender, sender.Addr())
		}

		email.Date = msg.Envelope.Date
		email.Subject = msg.Envelope.Subject

		emails = append(emails, &email)
	}

	return emails, nil
}

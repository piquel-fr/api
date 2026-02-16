package email

import (
	"context"
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/storage"
	"github.com/wneessen/go-mail"
)

const TrashFolder = "Trash" // TODO: save it in the mail account

type EmailHead struct {
	Id    uint32
	Flags []string
	To, Cc, Bcc,
	From, Sender []string
	Date    time.Time
	Subject string
}

type Email struct {
	Head EmailHead
	Body []byte
}

type Folder struct {
	Name        string `json:"name"`
	NumMessages int    `json:"num_messages"`
	NumUnread   int    `json:"num_unread"`
}

type EmailSendParams struct {
	Destinations []string `json:"destinations"`
	Subject      string   `json:"subject"`
	Content      string   `json:"content"`
}

type EmailService interface {
	// account stuff
	GetAccountByEmail(ctx context.Context, email string) (*repository.MailAccount, error)
	ListAccounts(ctx context.Context, userId int32) ([]*repository.MailAccount, error)
	CountAccounts(ctx context.Context, userId int32) (int64, error)
	AddAccount(ctx context.Context, params repository.AddEmailAccountParams) (int32, error)
	RemoveAccount(ctx context.Context, accountId int32) error
	GetAccountInfo(ctx context.Context, account *repository.MailAccount) (AccountInfo, error)

	// sharing
	AddShare(ctx context.Context, params repository.AddShareParams) error
	RemoveShare(ctx context.Context, userId, accountId int32) error
	GetAccountShares(ctx context.Context, account int32) ([]int32, error)

	SendEmail(from *repository.MailAccount, params EmailSendParams) error

	// folder management
	ListFolders(account *repository.MailAccount) ([]Folder, error)
	CreateFolder(account *repository.MailAccount, name string) error
	DeleteFolder(account *repository.MailAccount, name string) error
	RenameFolder(account *repository.MailAccount, name, newName string) error

	// email management
	GetFolderEmails(account *repository.MailAccount, folder string, offset, limit uint32) ([]*EmailHead, error)
	GetEmail(account *repository.MailAccount, folder string, id uint32) (Email, error)
	DeleteEmail(account *repository.MailAccount, folder string, id uint32) error
}

type realEmailService struct {
	imapAddr       string
	storageService storage.StorageService
}

func NewRealEmailService(storageService storage.StorageService) EmailService {
	addr := fmt.Sprintf("%s:%s", config.Envs.ImapHost, config.Envs.ImapPort)
	return &realEmailService{
		imapAddr:       addr,
		storageService: storageService,
	}
}

func (s *realEmailService) SendEmail(from *repository.MailAccount, params EmailSendParams) error {
	message := mail.NewMsg()
	if err := message.From(from.Email); err != nil {
		return fmt.Errorf("failed to add FROM address %s: %w", from.Email, err)
	}

	for _, to := range params.Destinations {
		if err := message.AddTo(to); err != nil {
			return fmt.Errorf("failed to add TO address %s: %w", to, err)
		}
	}

	message.Subject(params.Subject)
	message.SetBodyString(mail.TypeTextHTML, params.Content)

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

func (s *realEmailService) ListFolders(account *repository.MailAccount) ([]Folder, error) {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return nil, err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return nil, err
	}

	listCmd := client.List("", "*", &imap.ListOptions{ReturnStatus: &imap.StatusOptions{NumMessages: true, NumUnseen: true}})
	defer listCmd.Close()

	folders := []Folder{}
	for mailbox := listCmd.Next(); mailbox != nil; mailbox = listCmd.Next() {
		folders = append(folders, Folder{
			Name:        mailbox.Mailbox,
			NumMessages: int(*mailbox.Status.NumMessages),
			NumUnread:   int(*mailbox.Status.NumUnseen),
		})
	}

	return folders, nil
}

func (s *realEmailService) CreateFolder(account *repository.MailAccount, name string) error {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return err
	}

	return client.Create(name, nil).Wait()
}

func (s *realEmailService) DeleteFolder(account *repository.MailAccount, name string) error {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return err
	}

	// TODO: validation (can't delete Trash, INBOX, Sent)

	return client.Delete(name).Wait()
}

func (s *realEmailService) RenameFolder(account *repository.MailAccount, name, newName string) error {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return err
	}

	return client.Rename(name, newName, nil).Wait()
}

func (s *realEmailService) GetFolderEmails(account *repository.MailAccount, folder string, offset, limit uint32) ([]*EmailHead, error) {
	if limit > config.MaxLimit {
		limit = config.MaxLimit
	}

	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return nil, err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return nil, err
	}

	mailbox, err := client.Select(folder, nil).Wait()
	if err != nil {
		return nil, err
	}
	if mailbox.NumMessages == 0 {
		return nil, nil
	}

	start := offset + 1 // IMAP indeces start at 1
	stop := min(offset+limit, mailbox.NumMessages)

	seqSet := imap.SeqSet{{Start: start, Stop: stop}}
	fetchOptions := &imap.FetchOptions{
		Flags:    true,
		Envelope: true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierHeader, Peek: true},
		},
	}

	messages, err := client.Fetch(seqSet, fetchOptions).Collect()
	if err != nil {
		return nil, err
	}

	emails := []*EmailHead{}
	for _, msg := range messages {
		emails = append(emails, makeMailhead(msg))
	}

	return emails, nil
}

func (s *realEmailService) GetEmail(account *repository.MailAccount, folder string, id uint32) (Email, error) {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return Email{}, err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return Email{}, err
	}

	if _, err := client.Select(folder, nil).Wait(); err != nil {
		return Email{}, err
	}

	bodySection := &imap.FetchItemBodySection{Specifier: imap.PartSpecifierText}
	uidSet := imap.UIDSet{{Start: imap.UID(id), Stop: imap.UID(id)}}
	fetchOptions := &imap.FetchOptions{
		Flags:    true,
		Envelope: true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierHeader},
			bodySection,
		},
	}

	messages, err := client.Fetch(uidSet, fetchOptions).Collect()
	if err != nil {
		return Email{}, err
	}
	if len(messages) != 1 {
		return Email{}, fmt.Errorf("major error when fetching email %d from folder %s, %d emails were found", id, folder, len(messages))
	}

	message := messages[0]
	return Email{
		Head: *makeMailhead(message),
		Body: message.FindBodySection(bodySection),
	}, nil
}

func (s *realEmailService) DeleteEmail(account *repository.MailAccount, folder string, id uint32) error {
	client, err := imapclient.DialTLS(s.imapAddr, nil)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return err
	}

	if _, err := client.Select(folder, nil).Wait(); err != nil {
		return err
	}

	uidSet := imap.UIDSet{{Start: imap.UID(id), Stop: imap.UID(id)}}
	_, err = client.Move(uidSet, TrashFolder).Wait()
	return err
}

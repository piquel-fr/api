package repository

const (
	ResourceUser         string = "user"
	ResourceDocsInstance string = "docs_instance"
	ResourceMailAccount  string = "email_account"
)

func (profile *User) GetResourceName() string { return ResourceUser }
func (profile *User) GetOwner() int32         { return profile.ID }

func (docs *DocsInstance) GetResourceName() string { return ResourceDocsInstance }
func (docs *DocsInstance) GetOwner() int32         { return docs.OwnerId }

func (account *MailAccount) GetResourceName() string { return ResourceMailAccount }
func (account *MailAccount) GetOwner() int32         { return account.OwnerId }

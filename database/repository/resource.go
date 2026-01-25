package repository

const (
	ResourceUser         string = "user"
	ResourceMailAccount  string = "email_account"
)

func (profile *User) GetResourceName() string { return ResourceUser }
func (profile *User) GetOwner() int32         { return profile.ID }

func (account *MailAccount) GetResourceName() string { return ResourceMailAccount }
func (account *MailAccount) GetOwner() int32         { return account.OwnerId }

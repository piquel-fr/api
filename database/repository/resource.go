package repository

func (account *MailAccount) GetResourceName() string {
	return "email_account"
}

func (account *MailAccount) GetOwner() int32 {
	return account.OwnerId
}

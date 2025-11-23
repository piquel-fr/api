package repository

func (docs *DocsInstance) GetResourceName() string {
	return "docs_instance"
}

func (docs *DocsInstance) GetOwner() int32 {
	return docs.OwnerId
}

func (account *MailAccount) GetResourceName() string {
	return "email_account"
}

func (account *MailAccount) GetOwner() int32 {
	return account.OwnerId
}

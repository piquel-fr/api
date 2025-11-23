package email

type EmailService interface {
	SendEmail()
}

type realEmailService struct{}

func NewRealEmailService() *realEmailService {
	return &realEmailService{}
}

func (r *realEmailService) SendEmail() {}

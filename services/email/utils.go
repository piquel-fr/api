package email

import "github.com/emersion/go-imap/v2/imapclient"

func makeMailhead(msg *imapclient.FetchMessageBuffer) *EmailHead {
	email := &EmailHead{Id: uint32(msg.UID)}

	for _, flag := range msg.Flags {
		email.Flags = append(email.Flags, string(flag))
	}

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
	return email
}

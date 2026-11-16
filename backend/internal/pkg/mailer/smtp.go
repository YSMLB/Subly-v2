package mailer

import(
	"fmt"
	"gopkg.in/gomail.v2"
)

type SMTPMailer struct{
	dialer *gomail.Dialer
	from string
}

func NewSMTPMailer(port int, host, user, pass string) *SMTPMailer{
	return  &SMTPMailer{
		dialer: gomail.NewDialer(host, port, user, pass),
		from: user,
	}
}

func (s *SMTPMailer)SendVerificationCode(toEmail, code string) error{
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("Subly <%s>", s.from))
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Код подтверждения регистрации Subly")


	//письмо
	body := fmt.Sprintf(`
        <div style="background:#0E131E; color:#fff; padding:24px; font-family:sans-serif; border-radius:16px;">
            <h2 style="color:#00E5FF; margin-top:0;">Subly ⚡</h2>
            <p style="color:#A0AEC0;">Твой код для подтверждения регистрации:</p>
            <div style="background:#171D2B; padding:16px; border-radius:12px; text-align:center; margin:16px 0;">
                <span style="font-size:32px; letter-spacing:8px; font-weight:bold; color:#00E5FF;">%s</span>
            </div>
            <p style="color:#7E8696; font-size:12px;">Код действует 10 минут. Если вы не запрашивали регистрацию, проигнорируйте письмо.</p>
        </div>
    `, code)

	msg.SetBody("text/html", body)
	return s.dialer.DialAndSend(msg)
}
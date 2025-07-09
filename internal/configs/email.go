package configs

type EmailConfig struct {
	MailHost        string
	MailPort        int
	MailUsername    string
	MailPassword    string
	MailFromAddress string
	FrontVerifyUrl  string
}

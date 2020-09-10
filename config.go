package backapi

type Config struct {
	Database  CfgDatabase `json:"database"`
	HTTP      CfgHTTP     `json:"http"`
	Mailer    CfgMailer   `json:"mailer"`
	Stripe    CfgStripe   `json:"stripe"`
	DomainURL string      `json:"domain_url"`
}

type CfgDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type CfgHTTP struct {
	Port string `json:"port"`
}

type CfgMailer struct {
	FromMail         string `json:"from_mail"`
	FromUser         string `json:"from_user"`
	MailjetAPIKey    string `json:"mailjet_api_key"`
	MailjetAPISecret string `json:"mailjet_api_secret"`
}

type CfgStripe struct {
	SecretKey string `json:"secret_key"`
}

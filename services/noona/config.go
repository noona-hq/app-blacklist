package noona

type Config struct {
	BaseURL         string `default:"http://localhost:31140"`
	AppStoreURL     string `default:"http://localhost:31130/week#settings-apps"`
	ClientID        string `default:""`
	ClientSecret    string `default:""`
	AppBaseURL      string `default:"http://localhost:8080"`
	AppWebhookToken string `default:"very-secure-token-secret"`
}

package noona

type Config struct {
	BaseURL               string `default:"http://localhost:31140"`
	AooStoreURL           string `default:"http://localhost:31130/week#settings-apps"`
	ClientID              string `default:""`
	ClientSecret          string `default:""`
	BlacklistBaseURL      string `default:"http://localhost:8080"`
	BlacklistWebhookToken string `default:"very-secure-token-secret"`
}

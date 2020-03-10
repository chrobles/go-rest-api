package types

// Config : captures app config from env
type Config struct {
	CoinMarketCap struct {
		BaseURL string `envconfig:"CMC_API_BASEURL"`
		Key     string `envconfig:"CMC_API_KEY"`
		Local   bool   `envconfig:"USE_LOCAL"`
	}
	AzureBlob struct {
	}
	CosmosDB struct {
	}
}

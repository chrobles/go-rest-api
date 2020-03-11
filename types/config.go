package types

// Config : captures app config from env
type Config struct {
	AzureBlob     AzureBlobConfig
	CosmosDB      CosmosDbConfig
	CoinMarketCap CoinMarketCapConfig
}

// CoinMarketCapConfig : config for cmc api client
type CoinMarketCapConfig struct {
	BaseURL       string `envconfig:"CMC_API_BASEURL"`
	Key           string `envconfig:"CMC_API_KEY"`
	LocalDataPath string `envconfig:"CMC_LOCAL_DATA_PATH"`
	UseLocal      bool   `envconfig:"CMC_USE_LOCAL"`
}

// AzureBlobConfig : config for az blob client
type AzureBlobConfig struct {
}

// CosmosDbConfig : config for cosmosdb client
type CosmosDbConfig struct {
}

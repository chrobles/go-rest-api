package types

// Config : captures app config from env
type Config struct {
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

// CosmosDbConfig : config for cosmosdb client
type CosmosDbConfig struct {
	Connstr string `envconfig:"CDB_CONN_STR"`
}

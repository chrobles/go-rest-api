package types

// Config : captures app config from env
type Config struct {
	Root          string
	CosmosDB      CosmosDbConfig
	CoinMarketCap CoinMarketCapConfig
}

// CoinMarketCapConfig : config for cmc api client
type CoinMarketCapConfig struct {
	BaseURL  string
	Key      string
	UseLocal bool
}

// CosmosDbConfig : config for cosmosdb client
type CosmosDbConfig struct {
	Connstr   string
	UseCosmos bool
}

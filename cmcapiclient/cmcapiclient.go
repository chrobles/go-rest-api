package cmcapiclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/chrobles/go-rest-api/types"
)

// Client : client for interacting with coinmarketcap API
type Client struct {
	BaseURL string
	Key     string
	Local   bool
}

// Init : apply client configuration
func (client *Client) Init(cfg types.Config) error {
	client.BaseURL = cfg.CoinMarketCap.BaseURL
	client.Key = cfg.CoinMarketCap.Key
	client.Local = cfg.CoinMarketCap.Local

	return nil
}

// NewMarketRequest : generate a request for a range of data from coinmarketcap API
func (client *Client) NewMarketRequest(start string, limit string) *http.Request {
	req, err := http.NewRequest("GET", client.BaseURL, nil)
	if err != nil {
		log.Print(err)
		return nil
	}

	query := url.Values{}
	query.Add("start", start)
	query.Add("limit", limit)
	query.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", client.Key)
	req.URL.RawQuery = query.Encode()

	return req
}

// Get : do a request and return the results
func (client *Client) Get(req *http.Request) (*types.MarketListings, error) {
	var (
		res     *http.Response
		resdata *types.MarketListings
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	_ = json.NewDecoder(res.Body).Decode(resdata)

	return resdata, nil
}

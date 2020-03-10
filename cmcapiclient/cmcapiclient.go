package cmcapiclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chrobles/go-rest-api/types"
)

// Client : client for interacting with coinmarketcap API
type Client struct {
	BaseURL  string
	Key      string
	UseLocal bool
}

// Configure : apply client configuration
func (client *Client) Configure(cfg types.Config) error {
	client.BaseURL = cfg.CoinMarketCap.BaseURL
	client.Key = cfg.CoinMarketCap.Key
	client.UseLocal = cfg.CoinMarketCap.UseLocal

	if !client.UseLocal && client.Key == "" {
		return errors.New("env var CMC_API_KEY missing")
	}

	log.Print("serving market data from disk")

	return nil
}

// NewMarketRequest : generate a request for a range of data from coinmarketcap API
func (client *Client) NewMarketRequest(start int, limit int) (*http.Request, error) {
	var (
		err   error
		query url.Values
		req   *http.Request
	)

	req, err = http.NewRequest("GET", client.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	query = url.Values{}
	query.Add("start", strconv.Itoa(start))
	query.Add("limit", strconv.Itoa(limit))
	query.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", client.Key)
	req.URL.RawQuery = query.Encode()

	return req, nil
}

// Do : do a request for market listings from cmc api
func (client *Client) Do(req *http.Request) (*types.MarketListings, error) {
	var (
		err     error
		res     *http.Response
		mktdata *types.MarketListings
	)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	mktdata = new(types.MarketListings)
	err = json.NewDecoder(res.Body).Decode(mktdata)
	if err != nil {
		return nil, err
	}

	return mktdata, nil
}

// DoLocal : read and return market listings from disk
func (client *Client) DoLocal(limit int) (*types.MarketListings, error) {
	var (
		data    []byte
		err     error
		mktdata *types.MarketListings
	)

	data, err = ioutil.ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &mktdata)
	if err != nil {
		return nil, err
	}

	if limit > len(mktdata.Data) {
		limit = len(mktdata.Data) + 1
	}
	mktdata.Data = mktdata.Data[:limit]

	return mktdata, nil
}

// GetMarketListings : get market listing data from disk or from cmc api
func (client *Client) GetMarketListings(start int, limit int) (*types.MarketListings, error) {
	if limit > 5000 {
		return nil, errors.New("limit must be less than or equal to 5000")
	}

	if !client.UseLocal {
		req, err := client.NewMarketRequest(start, limit)
		if err != nil {
			return nil, err
		}
		return client.Do(req)
	}
	return client.DoLocal(limit)
}

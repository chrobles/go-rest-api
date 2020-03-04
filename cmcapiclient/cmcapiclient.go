package cmcapiclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Client : client for interacting with coinmarketcap API
type Client struct {
	Key     string
	Address string
}

// RangeData : resulting json from query against cm API
type RangeData struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			USD struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

// NewRangeRequest : generate a request for a range of data from coinmarketcap API
func (client *Client) NewRangeRequest(start string, limit string) *http.Request {
	req, err := http.NewRequest("GET", client.Address, nil)
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
func (client *Client) Get(req *http.Request) RangeData {
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var resData RangeData
	_ = json.NewDecoder(res.Body).Decode(&resData)

	return resData
}

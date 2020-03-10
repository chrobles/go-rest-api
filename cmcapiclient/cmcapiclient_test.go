package cmcapiclient_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/chrobles/go-rest-api/cmcapiclient"
	"github.com/chrobles/go-rest-api/types"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	// discard config logging during test
	log.SetOutput(ioutil.Discard)

	var (
		cfgvalid   []types.Config
		cfginvalid []types.Config
		cmcclient  cmcapiclient.Client
		err        error
	)

	cfgvalid = []types.Config{
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "",
				LocalDataPath: "test-data-path",
				UseLocal:      true,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "test-api-key",
				LocalDataPath: "test-data-path",
				UseLocal:      true,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "test-api-key",
				LocalDataPath: "",
				UseLocal:      false,
			},
		},
	}
	cfginvalid = []types.Config{
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "",
				LocalDataPath: "",
				UseLocal:      false,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "",
				LocalDataPath: "test-data-path",
				UseLocal:      false,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "",
				LocalDataPath: "",
				UseLocal:      true,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				BaseURL:       "",
				Key:           "test-api-key",
				LocalDataPath: "",
				UseLocal:      true,
			},
		},
	}

	for _, c := range cfgvalid {
		err = cmcclient.Configure(c)
		assert.Nil(t, err)
	}

	for _, c := range cfginvalid {
		err = cmcclient.Configure(c)
		assert.NotNil(t, err)
	}
}

func TestNewMarketRequest(t *testing.T) {
	type query struct {
		start int
		limit int
	}

	var (
		cmcclient   cmcapiclient.Client
		err         error
		expectedURL string
		req         *http.Request
		queries     []query
		testkey     string
	)

	const maxint = int(^uint(0) >> 1)
	const minint = -(maxint - 1)

	queries = []query{
		{start: 1, limit: 10},
		{start: 20, limit: 1000},
		{start: 0, limit: 0},
		{start: minint, limit: maxint},
	}

	testkey = "test-api-key"
	cmcclient.Key = testkey

	for _, q := range queries {
		req, err = cmcclient.NewMarketRequest(q.start, q.limit)
		expectedURL = fmt.Sprintf("convert=USD&limit=%d&start=%d", q.limit, q.start)
		assert.Nil(t, err)
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, expectedURL, req.URL.RawQuery)
		assert.Equal(t, "application/json", req.Header.Get("Accepts"))
		assert.Equal(t, testkey, req.Header.Get("X-Cmc_pro_api_key"))
	}
}

func TestDo(t *testing.T) {
	// var (
	// 	cmcclient  Client
	// 	err        error
	// 	testserver *httptest.Server
	// )

	// // Start a local HTTP server
	// testserver = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	// 	// Test request parameters
	// 	equals(t, req.URL.String(), "/some/path")
	// 	// Send response to be tested
	// 	rw.Write([]byte(`OK`))
	// }))
	// // Close the server when test finishes
	// defer server.Close()

	// // Use Client & URL from our local test server
	// api := API{server.Client(), server.URL}
	// body, err := api.DoStuff()

	// ok(t, err)
	// equals(t, []byte("OK"), body)
}

func TestDoLocal(t *testing.T) {
	var (
		cmcclient cmcapiclient.Client
		datapath  string
		err       error
		limits    []int
		mktdata   *types.MarketListings
		pwd       string
	)

	limits = []int{1, 100, 2000}

	pwd, err = os.Getwd()
	assert.Nil(t, err)

	datapath = strings.Split(pwd, "go-rest-api/")[0] + "go-rest-api/data.json"
	cmcclient.LocalDataPath = datapath

	for _, l := range limits {
		mktdata, err = cmcclient.DoLocal(l)
		assert.Nil(t, err)
		assert.NotNil(t, mktdata)
	}
}

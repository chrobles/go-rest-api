package cmcapiclient_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
				Key:      "",
				UseLocal: true,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				Key:      "test-api-key",
				UseLocal: true,
			},
		},
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				Key:      "test-api-key",
				UseLocal: false,
			},
		},
	}
	cfginvalid = []types.Config{
		{
			CoinMarketCap: types.CoinMarketCapConfig{
				Key:      "",
				UseLocal: false,
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

	datapath = pwd + "/../data.json"
	cmcclient.LocalDataPath = datapath

	for _, l := range limits {
		mktdata, err = cmcclient.DoLocal(l)
		assert.Nil(t, err)
		assert.NotNil(t, mktdata)
	}
}

package tokenrate

import (
	"fmt"
	"net/http"
)

type TokenUsdPriceResult struct {
	Time         string  `json:"time"`
	AssetIdBase  string  `json:"asset_id_base"`
	AssetIdQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

func (c *Client) TokenUsdPrice(token string) (*TokenUsdPriceResult, error) {
	result := &TokenUsdPriceResult{}
	resp, err := c.client.R().
		SetResult(result).
		Get(fmt.Sprintf("https://rest.coinapi.io/v1/exchangerate/%s/USD", token))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("TokenUsdPrice resp.StatusCode: %d", resp.StatusCode())
	}

	return result, nil
}

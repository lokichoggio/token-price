package model

type GetTokenUsdPriceData struct {
	Time  string  `json:"time"`
	Token string  `json:"token"`
	Price float64 `json:"price"`
}

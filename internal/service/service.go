package service

import (
	"time"

	"token-price/internal/config"
	"token-price/internal/memory"
	"token-price/internal/model"
	"token-price/pkg/log"
	tokenrate "token-price/sdk/token-rate"
)

func GetTokenUsdPrice(token string) (interface{}, error) {
	// try get from memory
	value, exist := memory.Get(token)
	if exist {
		tokenUsdPriceResult := value.(*tokenrate.TokenUsdPriceResult)

		log.Infof("GetTokenUsdPrice from memory, data: %+v", tokenUsdPriceResult)

		return model.GetTokenUsdPriceData{
			Time:  tokenUsdPriceResult.Time,
			Token: token,
			Price: tokenUsdPriceResult.Rate,
		}, nil
	}

	// get from http
	result, err := tokenrate.NewClient(config.GetConfig().ApiKey).TokenUsdPrice(token)
	if err != nil {
		return nil, err
	}

	log.Infof("GetTokenUsdPrice from http, data: %+v", result)

	// write memory
	memory.Set(token, result, time.Second*10)

	return model.GetTokenUsdPriceData{
		Time:  result.Time,
		Token: token,
		Price: result.Rate,
	}, nil
}

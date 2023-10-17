package tokenrate

import (
	"testing"
)

func TestTokenUsdPrice(t *testing.T) {
	result, err := NewClient("").TokenUsdPrice("BTC")

	if err != nil {
		t.Errorf("error: %s", err)
		return
	}

	t.Logf("result: %+v", result)
}

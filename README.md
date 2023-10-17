## Quickstart

```shell
git clone git@github.com:lokichoggio/token-price.git
cd token-price
go mod tidy
# notice: replace YOUR_API_KEY in /etc/dev.yaml
go run cmd/main.go
```

## Usage

get token usd price

```shell
curl http://127.0.0.1:8080/api/v1/get_token_usd_price\?token\=ETH
```

metrics addr: http://127.0.0.1:8080/metrics

pprof addr: http://127.0.0.1:8080/debug/pprof/




package irail

import (
	"github.com/go-resty/resty"
)

const ua = "irail-csv/1.0 (https://github.com/meyskens/irail-csv)"

func getClient() *resty.Request {
	client := resty.New()
	client.SetHeader("User-Agent", ua)
	client.SetHostURL("https://api.irail.be")
	client.SetQueryParam("format", "json")
	client.SetQueryParam("lang", "nl")

	return client.R()
}

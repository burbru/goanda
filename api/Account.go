package api

import (
	"fmt"
	"strings"

	"github.com/burbru/goanda/models"
)

// GetPricing fetches the prricing for a list of instruments
func (api *API) GetPricing(instruments []string) (*models.Prices, error) {
	instrumentsQstr := strings.Join(instruments, ",")
	data, err := SendRequest("GET", api.context.ApiURL+"/v3/accounts/"+api.context.Account+"/pricing?instruments="+instrumentsQstr, nil)
	if err != nil {
		return nil, err
	}
	prices, err := parsePrices(&data)
	return &prices, err
}

func (api *API) GetCandles(instrument string, num int, granularity string, priceComponent PriceComponent) (*models.Candles, error) {
	qStr := fmt.Sprintf("?price=%s&granularity=%s&count=%d", priceComponent, granularity, num)
	data, err := SendRequest("GET", api.context.ApiURL+"/v3/accounts/"+api.context.Account+"/instruments/"+instrument+"/candles"+qStr, nil)
	if err != nil {
		return nil, err
	}
	candles, err := parseCandles(&data)
	return &candles, err
}

func (api *API) CreateNewCandlesRequest() *CandlesRequest {
	return &CandlesRequest{
		api:               api,
		instrument:        "",
		count:             0,
		granularity:       S5,
		priceComponent:    PriceComponentAskBidMid,
		from:              "",
		to:                "",
		includeFirst:      false,
		includeLast:       false,
		dailyAlignment:    0,
		alignmentTimezone: "",
		weeklyAlignment:   "",
	}
}

func (api *API) GetCandlesWithRequest(candlesRequest CandlesRequest) (*models.Candles, error) {
	return candlesRequest.Send()
}

// SetCount sets the number of candles to fetch
func (r *CandlesRequest) SetCount(count int) *CandlesRequest {
	r.count = count
	return r
}

// SetFrom sets the from date
func (r *CandlesRequest) SetFrom(from string) *CandlesRequest {
	r.from = from
	return r
}

// SetTo sets the to date
func (r *CandlesRequest) SetTo(to string) *CandlesRequest {
	r.to = to
	return r
}

// SetIncludeFirst sets the includeFirst flag
func (r *CandlesRequest) SetIncludeFirst(includeFirst bool) *CandlesRequest {
	r.includeFirst = includeFirst
	return r
}

// SetIncludeLast sets the includeLast flag
func (r *CandlesRequest) SetIncludeLast(includeLast bool) *CandlesRequest {
	r.includeLast = includeLast
	return r
}

// SetDailyAlignment sets the dailyAlignment flag
func (r *CandlesRequest) SetDailyAlignment(dailyAlignment int) *CandlesRequest {
	r.dailyAlignment = dailyAlignment
	return r
}

// SetAlignmentTimezone sets the alignmentTimezone flag
func (r *CandlesRequest) SetAlignmentTimezone(alignmentTimezone string) *CandlesRequest {
	r.alignmentTimezone = alignmentTimezone
	return r
}

// SetWeeklyAlignment sets the weeklyAlignment flag
func (r *CandlesRequest) SetWeeklyAlignment(weeklyAlignment string) *CandlesRequest {
	r.weeklyAlignment = weeklyAlignment
	return r
}

// SetGranularity sets the granularity
func (r *CandlesRequest) SetGranularity(granularity Granularity) *CandlesRequest {
	r.granularity = granularity
	return r
}

// SetPriceComponent sets the price component
func (r *CandlesRequest) SetPriceComponent(priceComponent PriceComponent) *CandlesRequest {
	r.priceComponent = priceComponent
	return r
}

// SetInstrument sets the instrument
func (r *CandlesRequest) SetInstrument(instrument string) *CandlesRequest {
	r.instrument = instrument
	return r
}

// Send sends the request
func (r *CandlesRequest) Send() (*models.Candles, error) {
	qStr := fmt.Sprintf("?price=%s&granularity=%s", r.priceComponent, r.granularity)
	if r.from != "" {
		qStr += "&from=" + r.from
	}
	if r.to != "" {
		qStr += "&to=" + r.to
	}
	if r.includeFirst {
		qStr += "&includeFirst=true"
	}
	if r.includeLast {
		qStr += "&includeLast=true"
	}
	if r.dailyAlignment != 0 {
		qStr += fmt.Sprintf("&dailyAlignment=%d", r.dailyAlignment)
	}
	if r.alignmentTimezone != "" {
		qStr += "&alignmentTimezone=" + r.alignmentTimezone
	}
	if r.weeklyAlignment != "" {
		qStr += "&weeklyAlignment=" + r.weeklyAlignment
	}
	PrintWithColor("GET %s\n", Green, r.api.context.ApiURL+"/v3/accounts/"+r.api.context.Account+"/instruments/"+r.instrument+"/candles"+qStr)
	data, err := SendRequest("GET", r.api.context.ApiURL+"/v3/accounts/"+r.api.context.Account+"/instruments/"+r.instrument+"/candles"+qStr, nil)
	if err != nil {
		return nil, err
	}

	candles, err := parseCandles(&data)

	return &candles, err
}

// Get the list of instruments for the account
func (api *API) GetInstruments() (*models.Instruments, error) {
	data, err := SendRequest("GET", api.context.ApiURL+"/v3/accounts/"+api.context.Account+"/instruments", nil)
	if err != nil {
		return nil, err
	}
	instruments, err := parseInstruments(&data)
	return &instruments, err
}

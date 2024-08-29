package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Instrument struct {
	Name                        string         `json:"name"`
	Type                        InstrumentType `json:"type"`
	DisplayName                 string         `json:"displayName"`
	PipLocation                 int            `json:"pipLocation"`
	DisplayPrecision            int            `json:"displayPrecision"`
	TradeUnitsPrecision         int            `json:"tradeUnitsPrecision"`
	MinimumTradeSize            Float64String  `json:"minimumTradeSize"`
	MaximumTrailingStopDistance Float64String  `json:"maximumTrailingStopDistance"`
	MinimumTrailingStopDistance Float64String  `json:"minimumTrailingStopDistance"`
	MaximumPositionSize         Float64String  `json:"maximumPositionSize"`
	MaximumOrderUnits           Float64String  `json:"maximumOrderUnits"`
	MarginRate                  Float64String  `json:"marginRate"`
	GuaranteedStopLossOrderMode Mode           `json:"guaranteedStopLossOrderMode"`
	Tags                        []Tag          `json:"tags"`
	Financing                   Financing      `json:"financing"`
}

type Mode string

const (
	ModeDisabled Mode = "DISABLED"
	ModeAllowed  Mode = "ALLOWED"
	ModeRequired Mode = "REQUIRED"
)

type InstrumentType string

const (
	Currency InstrumentType = "CURRENCY"
	Cfd      InstrumentType = "CFD"
	Metal    InstrumentType = "METAL"
)

type Tag struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Financing struct {
	LongRate            Float64String  `json:"longRate"`
	ShortRate           Float64String  `json:"shortRate"`
	FinancingDaysOfWeek []FinancingDay `json:"financingDaysOfWeek"`
}

type DayOfWeek string

const (
	Monday    DayOfWeek = "MONDAY"
	Tuesday   DayOfWeek = "TUESDAY"
	Wednesday DayOfWeek = "WEDNESDAY"
	Thursday  DayOfWeek = "THURSDAY"
	Friday    DayOfWeek = "FRIDAY"
	Saturday  DayOfWeek = "SATURDAY"
	Sunday    DayOfWeek = "SUNDAY"
)

type FinancingDay struct {
	DayOfWeek   DayOfWeek `json:"dayOfWeek"`
	DaysCharged int       `json:"daysCharged"`
}

type Instruments struct {
	Instruments []Instrument `json:"instruments"`
}

type Float64String float64

func (f *Float64String) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("Float64String: %w", err)
	}
	*f = Float64String(value)
	return nil
}

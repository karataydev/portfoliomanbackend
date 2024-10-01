package investmentgrowth

import "time"

type GrowthDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type GrowthResult struct {
	WeekData       []GrowthDataPoint `json:"weekData"`
	MonthData      []GrowthDataPoint `json:"monthData"`
	ThreeMonthData []GrowthDataPoint `json:"threeMonthData"`
	YearData       []GrowthDataPoint `json:"yearData"`
}

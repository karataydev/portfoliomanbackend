package investmentgrowth

import (
	"fmt"
	"sort"
	"time"

	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/portfolio"
)

type Service struct {
	portfolioService *portfolio.Service
	assetService     *asset.Service
}

func NewService(portfolioService *portfolio.Service, assetService *asset.Service) *Service {
	return &Service{
		portfolioService: portfolioService,
		assetService:     assetService,
	}
}

func (s *Service) CalculateInvestmentGrowth(symbol string) (*GrowthResult, error) {
	// First, try to get a portfolio with this symbol
	portfolioInfo, err := s.portfolioService.GetPortfolioBySymbol(symbol)
	if err == nil {
		// If a portfolio is found, use the portfolio growth calculation
		return s.CalculatePortfolioInvestmentGrowth(portfolioInfo.Id)
	}

	// If not found as a portfolio, try to get an asset with this symbol
	assetInfo, err := s.assetService.GetAssetBySymbol(symbol)
	if err == nil {
		// If an asset is found, use the asset growth calculation
		return s.CalculateAssetInvestmentGrowth(assetInfo.Id)
	}

	// If neither a portfolio nor an asset is found, return an error
	return nil, fmt.Errorf("symbol %s not found as either portfolio or asset", symbol)
}

func (s *Service) CalculatePortfolioInvestmentGrowth(portfolioId int64) (*GrowthResult, error) {
	initialInvestment := 1000.0

	portfolio, err := s.portfolioService.GetPortfolioWithAllocations(portfolioId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	periods := map[string]time.Time{
		"week":       now.AddDate(0, 0, -7),
		"month":      now.AddDate(0, -1, 0),
		"threeMonth": now.AddDate(0, -3, 0),
		"year":       now.AddDate(-1, 0, 0),
	}

	result := &GrowthResult{}

	for period, startDate := range periods {
		growthData, err := s.calculateGrowthForPeriod(portfolio, initialInvestment, startDate, now, period)
		if err != nil {
			return nil, err
		}

		switch period {
		case "week":
			result.WeekData = growthData
		case "month":
			result.MonthData = growthData
		case "threeMonth":
			result.ThreeMonthData = growthData
		case "year":
			result.YearData = growthData
		}
	}

	return result, nil
}

func (s *Service) calculateGrowthForPeriod(portfolio *portfolio.PortfolioDTO, initialInvestment float64, startDate, endDate time.Time, period string) ([]GrowthDataPoint, error) {
	var allQuotes []asset.AssetQuote
	for _, allocation := range portfolio.Allocations {
		quotes, err := s.assetService.GetAssetQuotesForPeriod(allocation.Asset.Id, startDate, endDate)
		if err != nil {
			return nil, err
		}
		allQuotes = append(allQuotes, quotes...)
	}

	// Sort all quotes by timestamp
	sort.Slice(allQuotes, func(i, j int) bool {
		return allQuotes[i].QuoteTime.Before(allQuotes[j].QuoteTime)
	})

	var selectedQuotes []asset.AssetQuote
	var currentDay time.Time
	var dayQuotes []asset.AssetQuote

	for _, quote := range allQuotes {
		if !sameDay(currentDay, quote.QuoteTime) {
			if len(dayQuotes) > 0 {
				selectedQuotes = append(selectedQuotes, selectQuotesForPeriod(dayQuotes, period)...)
			}
			currentDay = quote.QuoteTime
			dayQuotes = []asset.AssetQuote{quote}
		} else {
			dayQuotes = append(dayQuotes, quote)
		}
	}
	if len(dayQuotes) > 0 {
		selectedQuotes = append(selectedQuotes, selectQuotesForPeriod(dayQuotes, period)...)
	}

	var growthData []GrowthDataPoint
	if len(selectedQuotes) > 0 {
		initialQuote := selectedQuotes[0].Quote
		for _, quote := range selectedQuotes {
			value := initialInvestment * (quote.Quote / initialQuote)
			growthData = append(growthData, GrowthDataPoint{
				Timestamp: quote.QuoteTime,
				Value:     value,
			})
		}
	}

	return growthData, nil
}

func (s *Service) CalculateAssetInvestmentGrowth(assetId int64) (*GrowthResult, error) {
	initialInvestment := 1000.0

	now := time.Now()
	periods := map[string]time.Time{
		"week":       now.AddDate(0, 0, -7),
		"month":      now.AddDate(0, -1, 0),
		"threeMonth": now.AddDate(0, -3, 0),
		"year":       now.AddDate(-1, 0, 0),
	}

	result := &GrowthResult{}

	for period, startDate := range periods {
		growthData, err := s.calculateAssetGrowthForPeriod(assetId, initialInvestment, startDate, now, period)
		if err != nil {
			return nil, err
		}

		switch period {
		case "week":
			result.WeekData = growthData
		case "month":
			result.MonthData = growthData
		case "threeMonth":
			result.ThreeMonthData = growthData
		case "year":
			result.YearData = growthData
		}
	}

	return result, nil
}

func (s *Service) calculateAssetGrowthForPeriod(assetId int64, initialInvestment float64, startDate, endDate time.Time, period string) ([]GrowthDataPoint, error) {
	quotes, err := s.assetService.GetAssetQuotesForPeriod(assetId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var selectedQuotes []asset.AssetQuote
	var currentDay time.Time
	var dayQuotes []asset.AssetQuote

	for _, quote := range quotes {
		if !sameDay(currentDay, quote.QuoteTime) {
			if len(dayQuotes) > 0 {
				selectedQuotes = append(selectedQuotes, selectQuotesForPeriod(dayQuotes, period)...)
			}
			currentDay = quote.QuoteTime
			dayQuotes = []asset.AssetQuote{quote}
		} else {
			dayQuotes = append(dayQuotes, quote)
		}
	}
	if len(dayQuotes) > 0 {
		selectedQuotes = append(selectedQuotes, selectQuotesForPeriod(dayQuotes, period)...)
	}

	var growthData []GrowthDataPoint
	if len(selectedQuotes) > 0 {
		initialQuote := selectedQuotes[0].Quote
		for _, quote := range selectedQuotes {
			value := initialInvestment * (quote.Quote / initialQuote)
			growthData = append(growthData, GrowthDataPoint{
				Timestamp: quote.QuoteTime,
				Value:     value,
			})
		}
	}

	return growthData, nil
}

func selectQuotesForPeriod(dayQuotes []asset.AssetQuote, period string) []asset.AssetQuote {
	if len(dayQuotes) == 0 {
		return nil
	}

	switch period {
	case "week":
		return dayQuotes // Return all data points for the week
	case "month":
		return selectNQuotes(dayQuotes, 4) // Select 4 data points per day
	case "threeMonth":
		return selectNQuotes(dayQuotes, 2) // Select 2 data points per day
	case "year":
		return []asset.AssetQuote{dayQuotes[len(dayQuotes)-1]} // Select last quote of the day
	default:
		return dayQuotes
	}
}

func selectNQuotes(quotes []asset.AssetQuote, n int) []asset.AssetQuote {
	if len(quotes) <= n {
		return quotes
	}
	step := float64(len(quotes)-1) / float64(n-1)
	var result []asset.AssetQuote
	for i := 0; i < n; i++ {
		index := int(float64(i) * step)
		result = append(result, quotes[index])
	}
	return result
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

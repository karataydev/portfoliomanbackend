package asset

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Service struct {
	repo          *Repository
	quoteReceiver <-chan AssetQuoteChanData
}

func NewService(repo *Repository, quoteReceiver <-chan AssetQuoteChanData) *Service {
	return &Service{
		repo:          repo,
		quoteReceiver: quoteReceiver,
	}
}

func (s *Service) GetAssets() ([]SimpleAssetDTO, error) {
	return s.repo.GetAssets()
}

func (s *Service) GetAsset(assetId int64) (*Asset, error) {
	return s.repo.GetAsset(assetId)
}

func (s *Service) GetAssetBySymbol(symbol string) (*Asset, error) {
	return s.repo.GetAssetBySymbol(symbol)
}

func (s *Service) GetAssetQuoteAtTime(assetId int64, t time.Time) (*AssetQuote, error) {
	return s.repo.GetAssetQuoteAtTime(assetId, t)
}

func (s *Service) GetAssetQuotesForPeriod(assetId int64, startTime, endTime time.Time) ([]AssetQuote, error) {
	return s.repo.GetAssetQuotesForPeriod(assetId, startTime, endTime)
}

func (s *Service) GetLatestQuote(assetId int64) (*AssetQuote, error) {
	return s.GetAssetQuoteAtTime(assetId, time.Now())
}

func (s *Service) SaveAssetQuote(assetQuoteData AssetQuoteChanData) error {

	assetQuote := AssetQuote{
		AssetId:   assetQuoteData.AssetId,
		Quote:     assetQuoteData.Quote,
		QuoteTime: assetQuoteData.QuoteTime,
	}

	return s.repo.SaveAssetQuote(assetQuote)
}

func (s *Service) AssetQuoteChanDataConsumer() {
	for quoteData := range s.quoteReceiver {
		err := s.SaveAssetQuote(quoteData)
		if err != nil {
			log.Error("Error consuming quote data.", quoteData, err)
		}
	}
}

func (s *Service) GetPreviousTradingDayQuote(assetId int64, currentTime time.Time) (*AssetQuote, error) {
	checkTime := currentTime.AddDate(0, 0, -1)

	for i := 0; i < 10; i++ { // Check up to 10 days back to be safe
		quote, err := s.GetAssetQuoteAtTime(assetId, checkTime)
		if err == nil {
			return quote, nil
		}

		// If no quote found, move to the previous day
		checkTime = checkTime.AddDate(0, 0, -1)
	}

	return nil, errors.New("no previous trading day quote found within the last 10 days")
}

func (s *Service) GetMarketOverview() ([]MarketGrowthListResponse, error) {
	assetSymbols := []string{"VOO", "AAPL", "GOOGL", "MSFT", "AMZN", "META"}
	assets, err := s.repo.GetAssetBySymbolList(assetSymbols)
	if err != nil {
		return nil, err
	}

	response := make([]MarketGrowthListResponse, 0, len(assets))
	for _, asset := range assets {
		// Get latest quote
		latestQuote, err := s.GetLatestQuote(asset.Id)
		if err != nil {
			return nil, err
		}

		// Get previous trading day quote
		previousTradingDayQuote, err := s.GetPreviousTradingDayQuote(asset.Id, latestQuote.QuoteTime)
		if err != nil {
			return nil, err
		}

		// Calculate daily change percentage for this asset
		assetChange := 0.0
		if previousTradingDayQuote.Quote != 0 {
			assetChange = ((latestQuote.Quote - previousTradingDayQuote.Quote) / previousTradingDayQuote.Quote) * 100
		}
		portfolioResponse := MarketGrowthListResponse{
			Id:     asset.Id,
			Symbol: asset.Symbol,
			Name:   asset.Name,
			Change: assetChange,
			Amount: latestQuote.Quote,
		}
		response = append(response, portfolioResponse)
	}

	return response, nil
}

package asset

import (
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

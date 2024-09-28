package asset

import (
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

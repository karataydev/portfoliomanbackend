package assetquotefeeder

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/param"
	"github.com/svarlamov/goyhfin"
)

type Service struct {
	assetService *asset.Service
	paramService *param.Service
	quoteChannel chan asset.AssetQuoteChanData
}

func NewService(assetService *asset.Service, paramService *param.Service, quoteChannel chan asset.AssetQuoteChanData) *Service {
	return &Service{
		assetService: assetService,
		paramService: paramService,
		quoteChannel: quoteChannel,
	}
}

func (s *Service) InsertInitialData() error {
	is, err := s.paramService.IsInitialDataInserted()
	if err != nil {
		log.Fatalf("could not run IsInitialDataInserted: %v", err)
	}
	if !is {
		if err := s.ScrapeAllAssets(goyhfin.OneYear, goyhfin.OneHour); err != nil {
			log.Fatalf("could not run scrape asssets: %v", err)
		}
		s.paramService.SetInitialDataInserted()
	} else {
		log.Info("Asset quote already inserted...")

	}
	return nil
}

func (s *Service) ScrapeAllAssets(rangeStr string, intervalStr string) error {
	assets, err := s.assetService.GetAssets()
	if err != nil {
		return err
	}

	for _, asset := range assets {
		err = s.ScrapeAsset(asset, rangeStr, intervalStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ScrapeAsset(asset asset.SimpleAssetDTO, rangeStr string, intervalStr string) error {
	resp, err := goyhfin.GetTickerData(asset.Symbol, rangeStr, intervalStr, false)
	if err != nil {
		return err
	}
	s.AssetResponseToChannel(asset.Id, resp)

	return nil
}

func (s *Service) AssetResponseToChannel(assetId int64, resp goyhfin.ChartQueryResponse) {
	for _, quote := range resp.Quotes {
		s.quoteChannel <- asset.AssetQuoteChanData{
			Symbol:    resp.Symbol,
			AssetId:   assetId,
			Quote:     quote.Close,
			QuoteTime: quote.ClosesAt,
		}
	}
}

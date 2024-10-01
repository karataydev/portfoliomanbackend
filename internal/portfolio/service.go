package portfolio

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/transaction"
)

type Service struct {
	repo               *Repository
	transactionService *transaction.Service
}

func NewService(repo *Repository, transactionService *transaction.Service) *Service {
	return &Service{
		repo:               repo,
		transactionService: transactionService,
	}
}

func (s *Service) GetPortfolio(portfolioId int64) (*PortfolioDTO, error) {
	return s.repo.GetPortfolio(portfolioId)
}

func (s *Service) GetAllocations(portfolioId int64) ([]AllocationDTO, error) {
	return s.repo.GetAllocations(portfolioId)
}

func (s *Service) GetPortfolioBySymbol(symbol string) (*Portfolio, error) {
    return s.repo.GetPortfolioBySymbol(symbol)
}

func (s *Service) GetPortfolioWithAllocations(portfolioId int64) (*PortfolioDTO, error) {
	portfolio, err := s.repo.GetPortfolioWithAllocations(portfolioId)
	if err != nil {
		return nil, err
	}

	var allocationIds []int64
	for _, allocation := range portfolio.Allocations {
		allocationIds = append(allocationIds, allocation.Id)
	}

	amountMap, err := s.transactionService.CalculateAmounts(allocationIds...)
	if err != nil {
		return nil, err
	}

	sumAmount := 0.0
	for _, val := range amountMap {
		sumAmount += val
	}

	for i := range portfolio.Allocations {
		amount := amountMap[portfolio.Allocations[i].Id]
		log.Info("Amount:", amount)
		portfolio.Allocations[i].Amount = amount

		if sumAmount != 0 {
			percentage := (amount / sumAmount) * 100
			log.Info("Percentage:", percentage)
			portfolio.Allocations[i].CurrentPercentage = percentage
		} else {
			log.Info("Percentage: 0 (sum amount is zero)")
			portfolio.Allocations[i].CurrentPercentage = 0
		}
	}

	log.Info(portfolio)

	return portfolio, nil
}

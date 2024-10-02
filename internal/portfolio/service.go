package portfolio

import (
	"errors"
	"fmt"

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
	var assetIds []int64
	for _, allocation := range portfolio.Allocations {
		allocationIds = append(allocationIds, allocation.Id)
		assetIds = append(assetIds, allocation.Asset.Id)
	}

	amountMap, err := s.transactionService.CalculateAmountsAndPL(allocationIds, assetIds)
	if err != nil {
		return nil, err
	}

	sumAmount := 0.0
	for _, val := range amountMap {
		sumAmount += val.CurrentAmount
	}

	for i := range portfolio.Allocations {
		amount := amountMap[portfolio.Allocations[i].Id]
		portfolio.Allocations[i].Amount = amount.CurrentAmount
		portfolio.Allocations[i].UnrealizedPL = amount.UnrealizedPL

		if sumAmount != 0 {
			percentage := (amount.CurrentAmount / sumAmount) * 100
			portfolio.Allocations[i].CurrentPercentage = percentage
		} else {
			portfolio.Allocations[i].CurrentPercentage = 0
		}
	}

	log.Info(portfolio)

	return portfolio, nil
}

func (s *Service) AddTransactionToPortfolio(request AddTransactionRequest) (*PortfolioDTO, error) {
	portfolio, err := s.GetPortfolioWithAllocations(request.PortfolioId)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	var allocationId int64 = -1
	for _, allocation := range portfolio.Allocations {
		if allocation.Asset.Symbol == request.Symbol {
			allocationId = allocation.Id
			break
		}
	}

	if allocationId == -1 {
		return nil, errors.New("symbol does not exist in portfolio allocations")
	}

	newTransaction := &transaction.Transaction{
		AllocationId: allocationId,
		Side:         request.Side,
		Quantity:     request.Quantity,
		Price:        request.AvgPrice,
	}

	_, err = s.transactionService.Save(newTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	return s.GetPortfolioWithAllocations(request.PortfolioId)
}

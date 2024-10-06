package portfolio

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/transaction"
)

type Service struct {
	repo               *Repository
	transactionService *transaction.Service
	assetService       *asset.Service
}

func NewService(repo *Repository, transactionService *transaction.Service, assetService *asset.Service) *Service {
	return &Service{
		repo:               repo,
		transactionService: transactionService,
		assetService:       assetService,
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

func (s *Service) GetPortfolioListByUser(userId int64) ([]PortfolioListResponse, error) {
	portfolios, err := s.repo.GetPortfolioByUserIdWithAllocations(userId)
	if err != nil {
		return nil, err
	}

	response := make([]PortfolioListResponse, 0, len(portfolios))

	for _, portfolio := range portfolios {
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
		sumPreviousDayAmount := 0.0
		for _, allocation := range portfolio.Allocations {
			amount := amountMap[allocation.Id]
			sumAmount += amount.CurrentAmount

			// Get latest quote
			latestQuote, err := s.assetService.GetLatestQuote(allocation.Asset.Id)
			if err != nil {
				return nil, err
			}

			// Get previous trading day quote
			previousTradingDayQuote, err := s.assetService.GetPreviousTradingDayQuote(allocation.Asset.Id, latestQuote.QuoteTime)
			if err != nil {
				return nil, err
			}

			// Calculate previous trading day amount
			previousDayAmount := (amount.CurrentAmount / latestQuote.Quote) * previousTradingDayQuote.Quote
			sumPreviousDayAmount += previousDayAmount
		}

		// Calculate daily change percentage
		dailyChange := 0.0
		if sumPreviousDayAmount != 0 {
			dailyChange = ((sumAmount - sumPreviousDayAmount) / sumPreviousDayAmount) * 100
		}

		portfolioResponse := PortfolioListResponse{
			Id:     portfolio.Id,
			Symbol: portfolio.Symbol,
			Name:   portfolio.Name,
			Change: dailyChange,
			Owner:  "",
			Amount: sumAmount,
		}

		response = append(response, portfolioResponse)
	}
	return response, nil
}

func (s *Service) FollowPortfolio(userID, portfolioID int64) error {
	// Check if the portfolio exists
	portfolio, err := s.repo.GetPortfolio(portfolioID)
	if err != nil {
		return err
	}
	if portfolio == nil {
		return errors.New("portfolio not found")
	}

	// Check if the user is trying to follow their own portfolio
	if portfolio.UserId == userID {
		return errors.New("cannot follow your own portfolio")
	}

	// Check if already following
	isFollowing, err := s.repo.IsFollowing(userID, portfolioID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following this portfolio")
	}

	return s.repo.FollowPortfolio(userID, portfolioID)
}

func (s *Service) UnfollowPortfolio(userID, portfolioID int64) error {
	// Check if actually following
	isFollowing, err := s.repo.IsFollowing(userID, portfolioID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return errors.New("not following this portfolio")
	}

	return s.repo.UnfollowPortfolio(userID, portfolioID)
}

func (s *Service) GetFollowedPortfolioList(userId int64) ([]PortfolioListResponse, error) {
	portfolios, err := s.repo.GetFollowedPortfoliosWithAllocations(userId)
	if err != nil {
		return nil, err
	}
	response := make([]PortfolioListResponse, 0, len(portfolios))

	for _, portfolio := range portfolios {
		var totalChange float64
		var totalPercentage float64

		for _, allocation := range portfolio.Allocations {
			// Get latest quote
			latestQuote, err := s.assetService.GetLatestQuote(allocation.Asset.Id)
			if err != nil {
				return nil, err
			}

			// Get previous trading day quote
			previousTradingDayQuote, err := s.assetService.GetPreviousTradingDayQuote(allocation.Asset.Id, latestQuote.QuoteTime)
			if err != nil {
				return nil, err
			}

			// Calculate daily change percentage for this asset
			assetChange := 0.0
			if previousTradingDayQuote.Quote != 0 {
				assetChange = ((latestQuote.Quote - previousTradingDayQuote.Quote) / previousTradingDayQuote.Quote) * 100
			}

			// Weight the change by the target percentage
			weightedChange := assetChange * (allocation.TargetPercentage / 100)
			totalChange += weightedChange
			totalPercentage += allocation.TargetPercentage
		}

		// Normalize the change if total percentage is not exactly 100%
		if totalPercentage != 0 {
			totalChange = (totalChange / totalPercentage) * 100
		}

		portfolioResponse := PortfolioListResponse{
			Id:     portfolio.Id,
			Symbol: portfolio.Symbol,
			Name:   portfolio.Name,
			Change: totalChange,
			Owner:  "", // Assuming you have this field in your portfolio struct
			Amount: 0,  // As per your request, we're not calculating the actual amount
		}

		response = append(response, portfolioResponse)
	}
	return response, nil
}

func (s *Service) GetFollowerCount(portfolioID int64) (int, error) {
	return s.repo.GetFollowerCount(portfolioID)
}

func (s *Service) IsFollowing(userID, portfolioID int64) (bool, error) {
	return s.repo.IsFollowing(userID, portfolioID)
}

func (s *Service) CreatePortfolio(req CreatePortfolioRequest) (*PortfolioDTO, error) {
	// Create the portfolio
	portfolio := &Portfolio{
		UserId:      req.UserId,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	createdPortfolio, err := s.repo.CreatePortfolio(portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}

	// Create allocations
	allocations := make([]Allocation, len(req.Allocations))
	for i, alloc := range req.Allocations {
		allocations[i] = Allocation{
			PortfolioId:      createdPortfolio.Id,
			AssetId:          alloc.AssetId,
			TargetPercentage: alloc.TargetPercentage,
		}
	}

	err = s.repo.CreateAllocations(allocations)
	if err != nil {
		return nil, fmt.Errorf("failed to create allocations: %w", err)
	}

	// Fetch the created portfolio with allocations
	return s.GetPortfolioWithAllocations(createdPortfolio.Id)
}

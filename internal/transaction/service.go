package transaction

import "github.com/karataydev/portfoliomanbackend/internal/asset"

type Service struct {
	repo         *Repository
	assetService *asset.Service
}

func NewService(repo *Repository, assetService *asset.Service) *Service {
	return &Service{
		repo:         repo,
		assetService: assetService,
	}
}

func (s *Service) Get(allocationIds ...int64) ([]Transaction, error) {
	return s.repo.Get(allocationIds...)
}

func (s *Service) Save(t *Transaction) (*Transaction, error) {
	return s.repo.Save(t)
}

func (s *Service) CalculateAmountsAndPL(allocationIds, assetIds []int64) (map[int64]AmountAndPLResult, error) {
	transactions, err := s.repo.Get(allocationIds...)
	if err != nil {
		return nil, err
	}

	// Create a map to associate allocation IDs with asset IDs
	allocationToAsset := make(map[int64]int64)
	for i, allocID := range allocationIds {
		allocationToAsset[allocID] = assetIds[i]
	}

	// Group transactions by allocation ID
	transactionsByAllocation := make(map[int64][]Transaction)
	for _, t := range transactions {
		transactionsByAllocation[t.AllocationId] = append(transactionsByAllocation[t.AllocationId], t)
	}

	resultMap := make(map[int64]AmountAndPLResult)
	for allocID, txs := range transactionsByAllocation {
		latestQuote, err := s.assetService.GetLatestQuote(allocationToAsset[allocID])
		if err != nil {
			return nil, err
		}
		quantity := 0.0
		totalCost := 0.0
		for _, t := range txs {
			if t.Side == Buy {
				quantity += t.Quantity
				totalCost += t.Quantity * t.Price
			} else {
				quantity -= t.Quantity
				totalCost -= t.Quantity * t.Price
			}
		}
		currentAmount := quantity * latestQuote.Quote
		unrealizedPL := currentAmount - totalCost

		resultMap[allocID] = AmountAndPLResult{
			CurrentAmount: currentAmount,
			UnrealizedPL:  unrealizedPL,
		}
	}

	return resultMap, nil
}

func (s *Service) CalculateAmounts(allocationIds ...int64) (map[int64]float64, error) {
	transactions, err := s.repo.Get(allocationIds...)
	if err != nil {
		return nil, err
	}

	// Group transactions by allocation ID
	transactionsByAllocation := make(map[int64][]Transaction)
	for _, t := range transactions {
		transactionsByAllocation[t.AllocationId] = append(transactionsByAllocation[t.AllocationId], t)
	}

	amountMap := make(map[int64]float64)
	for key, val := range transactionsByAllocation {
		amount := 0.0
		for _, t := range val {
			if t.Side == Buy {
				amount += t.Price * t.Quantity
			} else {
				amount -= t.Price * t.Quantity
			}
		}
		amountMap[key] = amount
	}

	return amountMap, nil
}

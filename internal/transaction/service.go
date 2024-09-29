package transaction

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(allocationIds ...int64) ([]Transaction, error) {
	return s.repo.Get(allocationIds...)
}

func (s *Service) Save(t *Transaction) (*Transaction, error) {
	return s.repo.Save(t)
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

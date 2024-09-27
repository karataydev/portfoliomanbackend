package portfolio

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetPortfolio(portfolioId int64) (*PortfolioDTO, error) {
    return s.repo.GetPortfolio(portfolioId)
}

func (s *Service) GetAllocations(portfolioId int64) ([]AllocationDTO, error) {
    return s.repo.GetAllocations(portfolioId)
}

func (s *Service) GetPortfolioWithAllocations(portfolioId int64) (*PortfolioDTO, error) {
    return s.repo.GetPortfolioWithAllocations(portfolioId)
}
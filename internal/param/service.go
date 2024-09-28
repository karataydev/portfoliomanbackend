package param

import "database/sql"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(key string) (*Param, error) {
	return s.repo.Get(key)
}

func (s *Service) Save(param *Param) error {
	return s.repo.Save(param)
}

func (s *Service) IsInitialDataInserted() (bool, error) {
	param, err := s.Get("INITIAL_1Y_QUOTE_INSERT")
	if err != nil {
		// If the error is due to row not found, return false without error
		if err == sql.ErrNoRows {
			return false, nil
		}
		// For any other error, return it
		return false, err
	}

	// Check if the value is "TRUE"
	return param.Value == "TRUE", nil
}

func (s *Service) SetInitialDataInserted() error {
	return s.Save(&Param{Key: "INITIAL_1Y_QUOTE_INSERT", Value: "TRUE"})
}
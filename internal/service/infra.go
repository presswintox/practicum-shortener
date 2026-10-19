package service

type DatabaseRepository interface {
	Ping() error
}

type InfraService struct {
	db DatabaseRepository
}

func NewInfraService(db DatabaseRepository) *InfraService {
	return &InfraService{db: db}
}

func (s *InfraService) Ping() error {
	return s.db.Ping()
}

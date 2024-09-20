package application

import (
	"log"

	"github.com/amoonguses1/my-grpc-server/internal/port"
)

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}
func (s *BankService) FindCurrentBalance(account string) float64 {
	bankAccount, err := s.db.GetBankAccountByAccountNumber(account)
	if err != nil {
		log.Println("Error on FindCurrentBalance")
	}

	return bankAccount.CurrentBalance
}

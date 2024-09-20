package port

import (
	"time"

	db "github.com/amoonguses1/my-grpc-server/internal/adaptor/database"
	"github.com/google/uuid"
)

type DummyDatabasePort interface {
	Save(data *db.DummyOrm) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (db.DummyOrm, error)
}

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string) (db.BankAccountOrm, error)
	CreateExchangeRate(r db.BankExchangeRateOrm) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (db.BankExchangeRateOrm, error)
	CreateTransaction(acct db.BankAccountOrm, t db.BankTransactionOrm) (uuid.UUID, error)
}

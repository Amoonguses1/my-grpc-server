package database

import (
	"log"
	"time"

	dbank "github.com/amoonguses1/my-grpc-server/internal/application/domain/bank"
	"github.com/google/uuid"
)

func (a *DatabaseAdaptor) GetBankAccountByAccountNumber(acct string) (BankAccountOrm, error) {
	var bankAccountOrm BankAccountOrm

	if err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error; err != nil {
		log.Printf("Cannot find bank account number :%v", err)
		return bankAccountOrm, err
	}

	return bankAccountOrm, nil
}

func (a *DatabaseAdaptor) CreateExchangeRate(r BankExchangeRateOrm) (uuid.UUID, error) {
	if err := a.db.Create(r).Error; err != nil {
		return uuid.Nil, err
	}

	return r.ExchangeRateUuid, nil
}

func (a *DatabaseAdaptor) GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (BankExchangeRateOrm, error) {
	var exchangeRateOrm BankExchangeRateOrm
	err := a.db.First(&exchangeRateOrm, "from_currency = ? "+"AND to_currency = ? "+
		" AND (? BETWEEN valid_from_timestamp and valid_to_timestamp)", fromCur, toCur, ts).Error

	return exchangeRateOrm, err
}

func (a *DatabaseAdaptor) CreateTransaction(acct BankAccountOrm, t BankTransactionOrm) (uuid.UUID, error) {
	tx := a.db.Begin()

	if err := tx.Create(t).Error; err != nil {
		tx.Rollback()

		return uuid.Nil, err
	}

	newAmount := t.Amount
	if t.TransactionType == dbank.TransactionTypeOut {
		newAmount = -1 * t.Amount
	}

	newAccountBalance := acct.CurrentBalance + newAmount
	if err := tx.Model(&acct).Updates(
		map[string]interface{}{
			"current_balance": newAccountBalance,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()

	return t.TransactionUuid, nil
}

func (a *DatabaseAdaptor) CreateTransfer(transfer BankTransferOrm) (uuid.UUID, error) {
	if err := a.db.Create(transfer).Error; err != nil {
		return uuid.Nil, err
	}

	return transfer.TransferUuid, nil
}
func (a *DatabaseAdaptor) CreateTransferTransactionPair(fromAccountOrm BankAccountOrm,
	toAccountOrm BankAccountOrm, fromTransactionOrm BankTransactionOrm, toTransactionOrm BankTransactionOrm) (bool, error) {
	tx := a.db.Begin()

	if err := tx.Create(fromTransactionOrm).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Create(toTransactionOrm).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// recalculate current balance
	fromAccountBalanceNew := fromAccountOrm.CurrentBalance - fromTransactionOrm.Amount
	if err := tx.Model(&fromAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": fromAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	toAccountBalanceNew := toAccountOrm.CurrentBalance + toTransactionOrm.Amount
	if err := tx.Model(&toAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": toAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
func (a *DatabaseAdaptor) UpdateTransferStatus(transfer BankTransferOrm, status bool) error {
	if err := a.db.Model(&transfer).Updates(
		map[string]interface{}{
			"transfer_success": status,
			"updated_at":       time.Now(),
		},
	).Error; err != nil {
		return err
	}

	return nil
}

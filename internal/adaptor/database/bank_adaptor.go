package database

import "log"

func (a *DatabaseAdaptor) GetBankAccountByAccountNumber(acct string) (BankAccountOrm, error) {
	var bankAccountOrm BankAccountOrm

	if err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error; err != nil {
		log.Printf("Cannot find bank account number :%v", err)
		return bankAccountOrm, err
	}

	return bankAccountOrm, nil
}

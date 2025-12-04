package main

type MemoryStorage struct {
	accounts map[string]*Account
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		accounts: make(map[string]*Account),
	}
}

func (ms *MemoryStorage) SaveAccount(account *Account) error {
	ms.accounts[account.ID] = account
	return nil
}

func (ms *MemoryStorage) LoadAccount(accountID string) (*Account, error) {
	account, exists := ms.accounts[accountID]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (ms *MemoryStorage) GetAllAccounts() ([]*Account, error) {
	accounts := make([]*Account, 0, len(ms.accounts))
	for _, account := range ms.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

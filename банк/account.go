package main

import (
	"fmt"
	"strings"
	"time"
)

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	a.History = append(a.History, Transaction{
		Type:      "deposit",
		Amount:    amount,
		Timestamp: time.Now(),
		Details:   "",
	})

	if a.Storage != nil {
		a.Storage.SaveAccount(a)
	}
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.Balance < amount {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	a.History = append(a.History, Transaction{
		Type:      "withdraw",
		Amount:    amount,
		Timestamp: time.Now(),
		Details:   "",
	})

	if a.Storage != nil {
		a.Storage.SaveAccount(a)
	}
	return nil
}

func (a *Account) Transfer(to *Account, amount float64) error {
	if a.ID == to.ID {
		return ErrSameAccountTransfer
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.Balance < amount {
		return ErrInsufficientFunds
	}

	a.Balance -= amount
	to.Balance += amount

	a.History = append(a.History, Transaction{
		Type:      "transfer_out",
		Amount:    amount,
		Timestamp: time.Now(),
		Details:   to.ID,
	})
	to.History = append(to.History, Transaction{
		Type:      "transfer_in",
		Amount:    amount,
		Timestamp: time.Now(),
		Details:   a.ID,
	})

	if a.Storage != nil {
		a.Storage.SaveAccount(a)
		a.Storage.SaveAccount(to)
	}
	return nil
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) GetStatement() string {
	if len(a.History) == 0 {
		return "История транзакций пуста"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Выписка по счету %s (%s):\n", a.ID, a.Owner))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for i, tx := range a.History {
		timestamp := tx.Timestamp.Format("2006-01-02 15:04:05")
		switch tx.Type {
		case "deposit":
			sb.WriteString(fmt.Sprintf("%d. [%s] ПОПОЛНЕНИЕ: +%.2f\n", i+1, timestamp, tx.Amount))
		case "withdraw":
			sb.WriteString(fmt.Sprintf("%d. [%s] СНЯТИЕ: -%.2f\n", i+1, timestamp, tx.Amount))
		case "transfer_out":
			sb.WriteString(fmt.Sprintf("%d. [%s] ПЕРЕВОД НА %s: -%.2f\n", i+1, timestamp, tx.Details, tx.Amount))
		case "transfer_in":
			sb.WriteString(fmt.Sprintf("%d. [%s] ПЕРЕВОД ОТ %s: +%.2f\n", i+1, timestamp, tx.Details, tx.Amount))
		}
	}

	sb.WriteString(strings.Repeat("-", 50) + "\n")
	sb.WriteString(fmt.Sprintf("Текущий баланс: %.2f", a.Balance))
	return sb.String()
}

func generateAccountID() string {
	return fmt.Sprintf("ACC%06d", time.Now().UnixNano()%1000000)
}

package main

import (
	"errors"
	"time"
)

var (
	ErrInsufficientFunds   = errors.New("недостаточно средств на счете")
	ErrInvalidAmount       = errors.New("некорректная сумма: должна быть положительной")
	ErrAccountNotFound     = errors.New("счет не найден")
	ErrSameAccountTransfer = errors.New("нельзя перевести средства на тот же счет")
)

type Transaction struct {
	Type      string
	Amount    float64
	Timestamp time.Time
	Details   string
}

type Account struct {
	ID      string
	Owner   string
	Balance float64
	Storage Storage
	History []Transaction
}

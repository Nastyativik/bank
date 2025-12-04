package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BankApp struct {
	storage Storage
	scanner *bufio.Scanner
}

func NewBankApp() *BankApp {
	return &BankApp{
		storage: NewMemoryStorage(),
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (app *BankApp) getInput(prompt string) string {
	fmt.Print(prompt)
	if app.scanner.Scan() {
		return strings.TrimSpace(app.scanner.Text())
	}
	return ""
}

func (app *BankApp) getFloatInput(prompt string) (float64, error) {
	for {
		input := app.getInput(prompt)
		if input == "" {
			return 0, errors.New("ввод отменен")
		}
		amount, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Ошибка: введите корректное число")
			continue
		}
		return amount, nil
	}
}

func (app *BankApp) createAccount() {
	owner := app.getInput("Введите имя владельца счета: ")
	if owner == "" {
		fmt.Println("Имя не может быть пустым")
		return
	}

	account := &Account{
		ID:      generateAccountID(),
		Owner:   owner,
		Balance: 0,
		Storage: app.storage,
		History: []Transaction{},
	}

	app.storage.SaveAccount(account)
	fmt.Println("Счет успешно создан!")
	fmt.Println("ID счета:", account.ID)
	fmt.Println("Владелец:", account.Owner)
}

func (app *BankApp) findAccount() *Account {
	accountID := app.getInput("Введите ID счета: ")
	if accountID == "" {
		return nil
	}

	account, err := app.storage.LoadAccount(accountID)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return nil
	}
	return account
}

func (app *BankApp) deposit() {
	account := app.findAccount()
	if account == nil {
		return
	}

	amount, err := app.getFloatInput("Введите сумму для пополнения: ")
	if err != nil {
		fmt.Println("Операция отменена")
		return
	}

	err = account.Deposit(amount)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return
	}

	fmt.Println("Счет успешно пополнен на", amount)
}

func (app *BankApp) withdraw() {
	account := app.findAccount()
	if account == nil {
		return
	}

	amount, err := app.getFloatInput("Введите сумму для снятия: ")
	if err != nil {
		fmt.Println("Операция отменена")
		return
	}

	err = account.Withdraw(amount)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return
	}

	fmt.Println("Со счета снято", amount)
}

func (app *BankApp) transfer() {
	fromAccount := app.findAccount()
	if fromAccount == nil {
		return
	}

	toAccountID := app.getInput("Введите ID счета получателя: ")
	if toAccountID == "" {
		fmt.Println("Операция отменена")
		return
	}

	toAccount, err := app.storage.LoadAccount(toAccountID)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return
	}

	amount, err := app.getFloatInput("Введите сумму для перевода: ")
	if err != nil {
		fmt.Println("Операция отменена")
		return
	}

	err = fromAccount.Transfer(toAccount, amount)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return
	}

	fmt.Println("Перевод успешно выполнен!", amount, "переведено на счет", toAccountID)
}

func (app *BankApp) checkBalance() {
	account := app.findAccount()
	if account == nil {
		return
	}

	fmt.Println("Текущий баланс счета", account.ID, ":", account.GetBalance())
}

func (app *BankApp) getStatement() {
	account := app.findAccount()
	if account == nil {
		return
	}

	fmt.Println(account.GetStatement())
}

func (app *BankApp) showAllAccounts() {
	accounts, err := app.storage.GetAllAccounts()
	if err != nil {
		fmt.Println("Ошибка при получении списка счетов:", err.Error())
		return
	}

	if len(accounts) == 0 {
		fmt.Println("Нет созданных счетов")
		return
	}

	fmt.Println("Список всех счетов:")
	fmt.Println(strings.Repeat("-", 40))
	for _, account := range accounts {
		fmt.Println("ID:", account.ID, "| Владелец:", account.Owner, "| Баланс:", account.Balance)
	}
	fmt.Println(strings.Repeat("-", 40))
}

func (app *BankApp) Run() {
	fmt.Println("Добро пожаловать в консольный банковский сервис!")

	for {
		fmt.Println("Выберите операцию:")
		fmt.Println("1. Создать счет")
		fmt.Println("2. Пополнить счет")
		fmt.Println("3. Снять средства")
		fmt.Println("4. Перевести средства")
		fmt.Println("5. Проверить баланс")
		fmt.Println("6. Получить выписку")
		fmt.Println("7. Показать все счета")
		fmt.Println("8. Выйти")

		choice := app.getInput("Ваш выбор: ")

		switch choice {
		case "1":
			app.createAccount()
		case "2":
			app.deposit()
		case "3":
			app.withdraw()
		case "4":
			app.transfer()
		case "5":
			app.checkBalance()
		case "6":
			app.getStatement()
		case "7":
			app.showAllAccounts()
		case "8":
			fmt.Println("Спасибо за использование нашего сервиса! До свидания!")
			return
		default:
			fmt.Println("Неверный выбор. Пожалуйста, введите число от 1 до 8.")
		}
	}
}

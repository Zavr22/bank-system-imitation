package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Account struct {
	IBAN    string
	Balance float64
	Active  bool
}

type PaymentSystem struct {
	accounts           map[string]*Account
	emissionAccount    string
	destructionAccount string
}

type PaymentSystemInterface interface {
	EmitMoney(amount float64)
	DestroyMoney(amount float64)
	OpenAccount(iban string)
	TransferMoney(fromIBAN, toIBAN string, amount float64) error
	TransferMoneyJSON(transferJSON string) error
	GetAccountsInfo() string
}

func NewPaymentSystem() PaymentSystemInterface {
	ps := &PaymentSystem{
		accounts:           make(map[string]*Account),
		emissionAccount:    "BY00EMIS00000000000000000000",
		destructionAccount: "BY99DEST00000000000000000000",
	}
	// Инициализация специальных счетов
	ps.accounts[ps.emissionAccount] = &Account{IBAN: ps.emissionAccount, Balance: 0, Active: true}
	ps.accounts[ps.destructionAccount] = &Account{IBAN: ps.destructionAccount, Balance: 0, Active: true}
	return ps
}

func (ps *PaymentSystem) EmitMoney(amount float64) {
	ps.accounts[ps.emissionAccount].Balance += amount
}

func (ps *PaymentSystem) DestroyMoney(amount float64) {
	ps.accounts[ps.destructionAccount].Balance += amount
}

func (ps *PaymentSystem) OpenAccount(iban string) {
	ps.accounts[iban] = &Account{
		IBAN:    iban,
		Balance: 0,
		Active:  true,
	}
}

func (ps *PaymentSystem) TransferMoney(fromIBAN, toIBAN string, amount float64) error {
	fromAccount, ok := ps.accounts[fromIBAN]
	if !ok || !fromAccount.Active || fromAccount.Balance < amount {
		return errors.New("invalid source account")
	}

	toAccount, ok := ps.accounts[toIBAN]
	if !ok || !toAccount.Active {
		return errors.New("invalid destination account")
	}
	if fromAccount.Balance < amount {
		return errors.New("not enough money for transfer")
	}
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	return nil
}

func (ps *PaymentSystem) TransferMoneyJSON(transferJSON string) error {
	var transferData struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	err := json.Unmarshal([]byte(transferJSON), &transferData)
	if err != nil {
		return err
	}

	return ps.TransferMoney(transferData.From, transferData.To, transferData.Amount)
}

func (ps *PaymentSystem) GetAccountsInfo() string {
	data, _ := json.Marshal(ps.accounts)
	return string(data)
}

func main() {
	ps := NewPaymentSystem()
	ps.OpenAccount("BY12345678901234567890123456")
	ps.OpenAccount("BY98765432109876543210987654")

	fmt.Println("Initial accounts info:", ps.GetAccountsInfo())

	ps.EmitMoney(1000)
	fmt.Println("After emission:", ps.GetAccountsInfo())

	transferJSON := `{"from":"BY00EMIS00000000000000000000","to":"BY12345678901234567890123456","amount":500}`
	ps.TransferMoneyJSON(transferJSON)
	fmt.Println("After transfer:", ps.GetAccountsInfo())
}

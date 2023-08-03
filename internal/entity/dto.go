package entity

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AccountDepositInput struct {
	Id     int
	Amount int
}

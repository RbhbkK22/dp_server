package models

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Fio   string `json:"fio"`
	Post  string `json:"post"`
	Pass  string `json:"-"`
}

type Purchase struct {
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	PurchaseDate string `json:"purchase_date"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Photo       *string `json:"photo"`
	Category    string  `json:"category"`
	Brand       string  `json:"brand"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}

type Client struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Order struct {
	Id            int     `json:"id"`
	ClientName    string  `json:"client_name"`
	Comment       string  `json:"comment"`
	Date          string  `json:"datatime"`
	ManagerName   string  `json:"manager_name"`
	CollectorName *string `json:"collector_name"`
	Status        string  `json:"status"`
}

type Worker struct {
	Name     string `json:"fio"`
	Position string `json:"post"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddWorker struct {
	Name     string `json:"fio"`
	IdPosition int    `json:"idposit"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ItemInfo struct {
	ItemId int
	Name    string
	Photo   string
	Quality int
}

type Brand struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Category struct{
	Id int `json:"id"`
	Name string `json:"name"`
}

type Position struct{
	Id int `json:"id"`
	Name string `json:"name"`
}

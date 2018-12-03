package model


type Good struct {
	Abiid        int
	Mainname     string
	Subtitle     string
	Brandid      string
	Brandname    string
	CategoryId   string
	Categoryname string
	Price        int
	Stock        string
	IntStock	 int
}

type GoodPriceInfo struct {
	Abiid		string
	Mainname	string
	Price		int
	Stock		string
	RealPrice	int
	Num         int
	Num2		int
}

type EmailAccount struct {
	Sender string
	Pwd string
}
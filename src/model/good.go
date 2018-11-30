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
	Abiid		int
	Mainname	string
	Price		int
	Stock		string
	RealPrice	int
	Num         int
} 
package types

type WaitterStruct struct {
	Name     string  `json:"Name"`
	Distance float64 `json:"Distance"`
}

type TableStruct struct {
	ID          int `json:"ID"`
	Status      int `json:"Status"`
	CleanStatus int `json:"CleanStatus"`
}

type OrderStatusStruct struct {
	TableID int    `json:"tableID"`
	Name    string `json:"Name"`
	Gender  int    `json:"Gender"`
	Phone   string `json:"Phone"`
	Aldult  int    `json:"Aldult"`
	Child   int    `json:"Child"`
	Date    string `json:"Date"`
	Time    string `json:"Time"`
	Remark  string `json:"Remark"`
}

type OrderStruct struct {
	Name   string `json:"Name"`
	Phone  string `json:"Phone"`
	Gender int    `json:"Gender"`
	Date   string `json:"Date"`
	Time   string `json:"Time"`
	Aldult int    `json:"Aldult"`
	Child  int    `json:"Child"`
	Table  []int  `json:"Table"`
	Remark string `json:"Remark"`
	Notify int    `json:"Notify"`
}

type CustomerStruct struct {
	Name   string `json:"Name"`
	Phone  string `json:"Phone"`
	Gender int    `json:"Gender"`
	Date   string `json:"Date"`
	Time   string `json:"Time"`
	Notify int    `json:"Notify"`
}

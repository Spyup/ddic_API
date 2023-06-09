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
	TableID        int    `json:"tableID"`
	OrderName      string `json:"Name"`
	OrderPhone     string `json:"Phone"`
	NumberOfPeople int    `json:"NumberOfPeople"`
	OrderDateTime  string `json:"DateTime"`
	Remark         string `json:"Remark"`
}

type OrderStruct struct {
	Name           string `json:"Name"`
	Phone          string `json:"Phone"`
	DateTime       string `json:"DateTime"`
	NumberOfPeople int    `json:"NumberOfPeople"`
	Table          []int  `json:"table"`
	Remark         string `json:"Remark"`
}

type SliceOrderStruct []OrderStruct

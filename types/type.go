package types

type WaitterStruct struct {
	Name     string  `json:"Name"`
	Distance float64 `json:"Distance"`
}

type TableStruct struct {
	ID          string `json:"ID"`
	Status      int    `json:"Status"`
	CleanStatus int    `json:"CleanStatus"`
}

type OrderStruct struct {
	TableID        string `json:"tableID"`
	OrderName      string `json:"Name"`
	OrderPhone     string `json:"Phone"`
	NumberOfPeople string `json:"NumberOfPeople"`
	OrderDateTime  string `json:"DateTime"`
	Remark         string `json:"Remark"`
}

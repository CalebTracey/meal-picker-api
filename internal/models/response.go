package models

type CuisineResponse struct {
	Cuisine *Cuisine
	Message Message
}

type AllCuisinesResponse struct {
	Cuisines []*Cuisine
	Message  Message
}

type Message struct {
	ErrorLog  []ErrorLog `json:"ErrorLog,omitempty"`
	HostName  string     `json:"HostName,omitempty"`
	Status    string     `json:"Status,omitempty"`
	TimeTaken string     `json:"TimeTaken,omitempty"`
	Count     int        `json:"Count,omitempty"`
}

type ErrorLog struct {
	Status    string `json:"Status,omitempty"`
	Trace     string `json:"Trace,omitempty"`
	RootCause string `json:"RootCause,omitempty"`
}

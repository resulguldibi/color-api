package entity

type ResponseStatus struct {
	IsSucccess bool   `json:"issuccess"`
	Message    string `json:"message"`
}

type Color struct {
	R          int  `json:"r"`
	G          int  `json:"g"`
	B          int  `json:"b"`
	IsSelected bool `json:"-"`
}

type IEntity interface {
	Do()
}

func (color Color) Do() {}

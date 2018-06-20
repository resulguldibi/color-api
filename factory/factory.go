package factory

import (
	"resulguldibi/color-api/entity"
)

var factoryList = make(map[string]IFactory)

func InitFactoryList() {
	factoryList["Color"] = ColorFactory{}
}

type IFactory interface {
	GetInstance() entity.IEntity
}

type ColorFactory struct {
}

func (colorFactory ColorFactory) GetInstance() entity.IEntity {
	return &entity.Color{}
}

func GetEntityInstance(name string) entity.IEntity {
	return factoryList[name].GetInstance()
}

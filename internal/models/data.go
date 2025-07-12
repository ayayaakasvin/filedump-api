package models

type Data map[string]any 

func NewData () Data {
	return make(Data)
}
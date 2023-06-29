package lib

import "reflect"


func DataBuilder() *CommonDataStruct{
	return NewCommonDataStruct()
}

type CommonDataStruct struct {
    Title string
    Data map[string]interface{}
}

func NewCommonDataStruct() *CommonDataStruct {
	return &CommonDataStruct{Data: make(map[string]interface{})}
}
func(this *CommonDataStruct) SetTitle(title string) *CommonDataStruct{
	this.Title=title
	return this
}
func(this *CommonDataStruct) SetData(key string,value interface{}) *CommonDataStruct{
	this.Data[key]=value
	return this
}
func(this *CommonDataStruct) ToMap() (m map[string]interface{})  {
	m=make(map[string]interface{})
	elem := reflect.ValueOf(this).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return
}

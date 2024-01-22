package models

import (
	"reflect"
)

func init() {
	ModelList = append(ModelList, &Auth{})
	ModelList = append(ModelList, &Proxy{})
	ModelList = append(ModelList, &Tunnel{})
	ModelList = append(ModelList, &CloudProvider{})
	ModelList = append(ModelList, &SystemConfig{})
}

var ModelList = make([]interface{}, 0)

// ToApi 标准转化 API 方法
func ToApi(src interface{}, dst interface{}, extras ...func(api interface{})) interface{} {
	copyReflect(src, dst)
	// 自定义扩展
	for _, ex := range extras {
		ex(dst)
	}
	return dst
}

// ToModel 标准转化 Model 方法
func ToModel(src interface{}, dst interface{}) {
	copyReflect(src, dst)
}

func AutoFull(v interface{}) {
	val := reflect.ValueOf(v).Elem() // 获取指向结构体的反射值对象
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// 检查字段是否为指针并且为 nil
		if field.Kind() == reflect.Ptr && field.IsNil() {

			fieldType := field.Type().Elem() // 获取指针指向的类型
			newField := reflect.New(fieldType)

			// 根据字段类型创建新的实例
			switch fieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				newField.Elem().SetInt(0) // 设置为 0
			case reflect.Float32, reflect.Float64:
				newField.Elem().SetFloat(0.0) // 设置为 0.0
			case reflect.String:
				newField.Elem().SetString("") // 设置为空字符串
				// 可以根据需要添加更多类型
			}

			field.Set(newField)
		}
	}
}

// copyReflect 通过反射将相同的字段映射到对应的结构体中，减少重复的赋值动作
func copyReflect(src interface{}, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcType := srcVal.Type().Field(i)

		// 如果字段是匿名的结构体（可能是嵌入的结构体），则递归处理
		if srcType.Anonymous && srcField.Kind() == reflect.Struct {
			copyReflect(srcField.Interface(), dstVal.Addr().Interface())
			continue
		}

		// 查找目标结构中相同名称的字段
		dstField := dstVal.FieldByName(srcType.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue // 目标中没有这个字段，或者该字段不能被设置
		}

		// 确保源和目标字段类型相同
		if dstField.Type() == srcField.Type() {
			dstField.Set(srcField)
		}
	}
}

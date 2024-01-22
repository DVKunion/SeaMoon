package models

func init() {
	ModelList = append(ModelList, &Auth{})
	ModelList = append(ModelList, &Proxy{})
	ModelList = append(ModelList, &Tunnel{})
	ModelList = append(ModelList, &SystemConfig{})
	ModelList = append(ModelList, &CloudProvider{})
}

var ModelList = make([]interface{}, 0)

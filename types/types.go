package types

type Languages struct {
	En Language
	Sv Language
}

type Language struct {
	Key  string
	Name string
	Data map[string]interface{}
}


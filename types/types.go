package types

type BaseTranslation interface {
	Key() string
}

type Translation struct {
	key string
	En  string
	Sv  string
}

func (t Translation) Key() string {
	return t.key
}

func NewTranslation(key, en, sv string) Translation {
	return Translation{
		key: key,
		En:  en,
		Sv:  sv,
	}
}

// A group of translations or nested groups.
type TranslationGroup struct {
	key string

	// A list of translations or nested groups.
	Data []BaseTranslation
}

func (t TranslationGroup) Key() string {
	return t.key
}

func NewTranslationGroup(key string, data []BaseTranslation) TranslationGroup {
	return TranslationGroup{
		key:  key,
		Data: data,
	}
}

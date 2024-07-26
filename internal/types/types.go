package types

type JsonTranslations struct {
	En map[string]interface{}
	Sv map[string]interface{}
}

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

// A response to/reply from the Firebase Remote Config API.
//
// Parameters is the only important field in this application,
// so the other fields' types are not defined.
type FirebaseRemoteConfig struct {
	Conditions      []interface{}                            `json:"conditions"`
	Parameters      map[string]*FirebaseRemoteConfigParameter `json:"parameters"`
	Version         interface{}                              `json:"version"`
	ParameterGroups []interface{}                            `json:"parameterGroups"`
}

type FirebaseRemoteConfigParameter struct {
	DefaultValue      *FirebaseRemoteConfigValue            `json:"defaultValue"`
	ConditionalValues map[string]*FirebaseRemoteConfigValue `json:"conditionalValues"`
	Description       string                               `json:"description"`
	ValueType         string                               `json:"valueType"`
}

type FirebaseRemoteConfigValue struct {
	Value                *string      `json:"value"`
	UseInAppDefault      *bool        `json:"useInAppDefault"`
	PersonalizationValue *interface{} `json:"personalizationValue"`
	RolloutValue         *interface{} `json:"rolloutValue"`
}

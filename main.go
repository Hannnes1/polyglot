package main

import (
	"encoding/json"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Hannnes1/polyglot/types"
	"github.com/Hannnes1/polyglot/view"
)

func handleGetTranslations(c echo.Context) error {
	jsonDataEn := `{
        "confirm": "Confirm",
        "cancel": "Cancel",
        "example": {
            "title": "Example",
            "content": "This is an example content.",
			"subExample": {
				"title": "Subexample",
				"aTitle": "Subexample",
				"bTitle": "Subexample"
			}
        }
    }`

	jsonDataSv := `{
		"confirm": "Bekräfta",
		"cancel": "Avbryt",
		"example": {
			"title": "Exempel",
			"content": "Detta är ett exempel innehåll.",
			"subExample": {
				"title": "Underexempel",
				"aTitle": "Underexempel"
			}
		}
	}`

	var dataEn map[string]interface{}
	var dataSv map[string]interface{}

	errEn := json.Unmarshal([]byte(jsonDataEn), &dataEn)
	if errEn != nil {
		return errEn
	}
	errSv := json.Unmarshal([]byte(jsonDataSv), &dataSv)
	if errSv != nil {
		return errSv
	}

	t := buildTranslations(dataEn, dataSv)

	component := view.Index(t)

	component.Render(c.Request().Context(), c.Response().Writer)

	return nil
}

func buildTranslations(dataEn map[string]interface{}, dataSv map[string]interface{}) []types.BaseTranslation {
	var translations []types.BaseTranslation

	for key, value := range dataEn {
		switch v := value.(type) {
		case string:
			svString, _ := dataSv[key].(string)

			translations = append(translations, types.NewTranslation(
				key,
				v,
				svString,
			))

			// TODO: Sort alphabetically.
		case map[string]interface{}:
			translations = append(translations, types.NewTranslationGroup(
				key,
				buildTranslations(v, dataSv[key].(map[string]interface{})),
			))
		}
	}

	sort.Slice(translations, func(i, j int) bool {
		return translations[i].Key() < translations[j].Key()
	})

	return translations
}

func handlePostTranslations(c echo.Context) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		return err
	}

	svMap := jsonMap["sv"].(map[string]interface{})
	enMap := jsonMap["en"].(map[string]interface{})

	svString, svErr := json.MarshalIndent(svMap, "", "    ")
	if svErr != nil {
		return svErr
	}

	enString, enErr := json.MarshalIndent(enMap, "", "    ")
	if enErr != nil {
		return enErr
	}

	println(string(svString))
	println(string(enString))

  component := view.Confirmation()

	component.Render(c.Request().Context(), c.Response().Writer)

	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.Static("/assets", "assets")

	e.GET("/", handleGetTranslations)

	e.POST("/", handlePostTranslations)

	e.Logger.Fatal(e.Start(":8080"))
}

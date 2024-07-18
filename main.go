package main

import (
	"encoding/json"

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
				"title": "Subexample"
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
				"title": "Underexempel"
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

	l := types.Languages{
		En: types.Language{Key: "en", Name: "English", Data: dataEn},
		Sv: types.Language{Key: "sv", Name: "Swedish", Data: dataSv},
	}

  component := view.Index(l)

  component.Render(c.Request().Context(), c.Response().Writer)

  return nil
}

func main() {
  e := echo.New()

  e.Use(middleware.Logger())

  e.Static("/assets", "assets")

  e.GET("/", handleGetTranslations)
  
  e.Logger.Fatal(e.Start(":8080"))
}

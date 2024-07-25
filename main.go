package main

import (
	"encoding/json"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
  "github.com/joho/godotenv"

	"github.com/Hannnes1/polyglot/internal/repository"
	"github.com/Hannnes1/polyglot/internal/types"
	"github.com/Hannnes1/polyglot/internal/view"
)

func handleGetTranslations(c echo.Context, repo repository.BaseRepository) error {

	jsonData, err := repo.GetTranslations()

	if err != nil {
		return err
	}

	t := buildTranslations(jsonData.En, jsonData.Sv)

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

func handlePostTranslations(c echo.Context, repo repository.BaseRepository) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		return err
	}

	svJson := jsonMap["sv"].(map[string]interface{})
	enJson := jsonMap["en"].(map[string]interface{})

	repo.SaveTranslation(&types.JsonTranslations{
		En: enJson,
		Sv: svJson,
	})

	component := view.Confirmation()

	component.Render(c.Request().Context(), c.Response().Writer)

	return nil
}

func main() {
  err := godotenv.Load()
  if err != nil {
    panic(err)
  }

	repo, err := repository.NewFirebaseRepository()
  if err != nil {
    panic(err)
  }

	e := echo.New()

	e.Use(middleware.Logger())

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		return handleGetTranslations(c, repo)
	})

	e.POST("/", func(c echo.Context) error {
		return handlePostTranslations(c, repo)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

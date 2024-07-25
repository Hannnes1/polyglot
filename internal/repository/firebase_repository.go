package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/Hannnes1/polyglot/internal/types"
	"golang.org/x/oauth2/google"
)

var client *http.Client
var remoteConfigUrl string

type FirebaseRepository struct {
}

func NewFirebaseRepository() (*FirebaseRepository, error) {
	ctx := context.Background()

  remoteConfigUrl = "https://firebaseremoteconfig.googleapis.com/v1/projects/" + os.Getenv("GCP_PROJECT") + "/remoteConfig"

	c, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/firebase.remoteconfig")
	if err != nil {
		return nil, err
	}

	client = c

	return &FirebaseRepository{}, nil
}

func (f FirebaseRepository) GetTranslations() (*types.JsonTranslations, error) {
	res, err := client.Get(remoteConfigUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

  if res.StatusCode != 200 {
    io.Copy(os.Stdout, res.Body)

    return nil, errors.New("Failed to fetch translations")
  }

	body := types.FirebaseRemoteConfigResponse{}

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	var enMap map[string]interface{}
	var svMap map[string]interface{}

	err = json.Unmarshal([]byte(body.Parameters["translations_en"].DefaultValue.Value), &enMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body.Parameters["translations_sv"].DefaultValue.Value), &svMap)
	if err != nil {
		return nil, err
	}

	t := types.JsonTranslations{
		En: enMap,
		Sv: svMap,
	}

	return &t, nil
}

func (f FirebaseRepository) SaveTranslation(translation *types.JsonTranslations) error {
	en, err := json.Marshal(translation.En)
	if err != nil {
		return err
	}

	sv, err := json.Marshal(translation.Sv)
	if err != nil {
		return err
	}

	enString := string(en)
	svString := string(sv)

	enValue := types.FirebaseRemoteConfigValue{
		Value: enString,
	}

	svValue := types.FirebaseRemoteConfigValue{
		Value: svString,
	}

	enParameter := types.FirebaseRemoteConfigParameter{
		DefaultValue: enValue,
	}

	svParameter := types.FirebaseRemoteConfigParameter{
		DefaultValue: svValue,
	}

	parameters := map[string]types.FirebaseRemoteConfigParameter{
		"translations_en": enParameter,
		"translations_sv": svParameter,
	}

	body := types.FirebaseRemoteConfigResponse{
		Parameters: parameters,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", remoteConfigUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

  req.Header.Set("If-Match", "*")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

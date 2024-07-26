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

// Get the latest Firebase Remote Config, and associated ETag.
func (f FirebaseRepository) getConfig() (*types.FirebaseRemoteConfig, string, error) {
	req, err := http.NewRequest("GET", remoteConfigUrl, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("x-goog-user-project", os.Getenv("GCP_PROJECT"))

	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		println(res.StatusCode)
		io.Copy(os.Stdout, res.Body)

		return nil, "", errors.New("Failed to fetch translations")
	}

	body := types.FirebaseRemoteConfig{}

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, "", err
	}

	return &body, res.Header.Get("ETag"), nil
}

func (f FirebaseRepository) GetTranslations() (*types.JsonTranslations, error) {
	config, _, err := f.getConfig()
	if err != nil {
		return nil, err
	}

	var enMap map[string]interface{}
	var svMap map[string]interface{}

	err = json.Unmarshal([]byte(*config.Parameters["translations_en"].DefaultValue.Value), &enMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(*config.Parameters["translations_sv"].DefaultValue.Value), &svMap)
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

	config, eTag, err := f.getConfig()
	if err != nil {
		return err
	}

	config.Parameters["translations_en"].DefaultValue.Value = &enString
	config.Parameters["translations_sv"].DefaultValue.Value = &svString

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", remoteConfigUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("If-Match", eTag)
	req.Header.Set("x-goog-user-project", os.Getenv("GCP_PROJECT"))

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		println(res.StatusCode)
		io.Copy(os.Stdout, res.Body)

		return errors.New("Failed to save translations")
	}

	return nil
}

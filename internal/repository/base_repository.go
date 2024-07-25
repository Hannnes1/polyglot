package repository

import "github.com/Hannnes1/polyglot/internal/types"

type BaseRepository interface {
	GetTranslations() (*types.JsonTranslations, error)
	SaveTranslation(*types.JsonTranslations) error
}

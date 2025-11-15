package store

import (
    "errors"
    "be_test/internal/model"
)

var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid")
)

type ItemStore interface {
    Create(name string) (model.Item, error)
    Get(id string) (model.Item, error)
    List() ([]model.Item, error)
    Update(id string, name *string, done *bool) (model.Item, error)
    Delete(id string) error
}
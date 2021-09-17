package uuid

import (
	"github.com/google/uuid"
)

type Google struct {
}

func NewGoogleUUID() *Google {
	return &Google{}
}

func (g *Google) Generate() string {
	return uuid.New().String()
}

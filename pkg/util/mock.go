package util

import (
	"testing"
	"words/domain/repository"
)

func TestMain(m *testing.M) {
	repository.InitMongo()
	m.Run()
}

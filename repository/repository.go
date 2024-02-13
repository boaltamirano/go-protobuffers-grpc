package repository

import (
	"context"

	"github.com/boaltamirano/go-protobuffers-grpc/models"
)

type Repository interface {
	//************************** STUDENT ******************************************//
	GetStudent(ctx context.Context, id string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error

	//************************** TEST MODELS ******************************************//
	GetTest(ctx context.Context, id string) (*models.Test, error)
	SetTest(ctx context.Context, test *models.Test) error

	//************************** QUESTION MODELS ******************************************//
	SetQuestion(ctx context.Context, question *models.Question) error
}

// ************************** REPOSITORY IMPLEMENTATION ******************************************//
var implementation Repository

// ************************** REPOSITORY CONSTRUCTOR ******************************************//
func SetRepository(repository Repository) {
	implementation = repository
}

// ************************** STUDENT METHODS ******************************************//
func SetStudent(ctx context.Context, student *models.Student) error {
	return implementation.SetStudent(ctx, student)
}

func GetStudent(ctx context.Context, id string) (*models.Student, error) {
	return implementation.GetStudent(ctx, id)
}

// ************************** TEST METHODS ******************************************//
func GetTest(ctx context.Context, id string) (*models.Test, error) {
	return implementation.GetTest(ctx, id)
}

func SetTest(ctx context.Context, test *models.Test) error {
	return implementation.SetTest(ctx, test)
}

// ************************** QUESTION METHODS ******************************************//
func SetQuestion(ctx context.Context, question *models.Question) error {
	return implementation.SetQuestion(ctx, question)
}

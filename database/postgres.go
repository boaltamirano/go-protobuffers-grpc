package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/boaltamirano/go-protobuffers-grpc/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

// Constructor
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

// ************************** STUDENT ******************************************//
func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)", student.Id, student.Name, student.Age)
	return err
}

func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)

	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var student = models.Student{}

	for rows.Next() {
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
		return &student, nil
	}
	return &student, nil
}

// ************************** TEST MODELS ******************************************//
func (repo *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO tests (id, name) VALUES ($1, $2)", test.Id, test.Name)
	return err
}

func (repo *PostgresRepository) GetTest(ctx context.Context, id string) (*models.Test, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name FROM tests WHERE id = $1", id)

	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var testModel = models.Test{}

	for rows.Next() {
		err := rows.Scan(&testModel.Id, &testModel.Name)
		if err != nil {
			return nil, err
		}
		return &testModel, nil
	}
	return &testModel, nil
}

// ************************** QUESTION MODELS ******************************************//
func (repo *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO questions (id, question, answer, test_id) VALUES ($1, $2, $3, $4)", question.Id, question.Question, question.Answer, question.TestId)
	return err
}

// ************************** ENROLLMENT MODELS ******************************************//
func (repo *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testId)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var students []*models.Student
	for rows.Next() {
		var student = models.Student{}
		if err = rows.Scan(&student.Id, &student.Name, &student.Age); err == nil {
			students = append(students, &student)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil

}

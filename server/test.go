// ESTE ARCHIVO ES EL SERVIDOR DEL MODELO TEST

package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/boaltamirano/go-protobuffers-grpc/models"
	"github.com/boaltamirano/go-protobuffers-grpc/repository"
	"github.com/boaltamirano/go-protobuffers-grpc/studentpb"
	"github.com/boaltamirano/go-protobuffers-grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

// Constructor:
// // parametros repository.Repository
func NewTestServer(repo repository.Repository) *TestServer {
	// retorna una instancia de la struct TestServer, y le pasamos el repositorio que recibe como parametro
	return &TestServer{repo: repo}
}

// Implementacion de SetTest y GetTest a nivel de este servidor, estos metodos debes recibir y retornar lo que se definio en test.proto

// El GetTest retorna *testpb.Test que compilo el archivo test.proto
func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())

	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

// QUESTIONS
func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		msg, err := stream.Recv() // esta funcion stream.Recv() se bloquea hasta que llege los mensajes del cliente
		if err == io.EOF {        //EOF sirve para identificar un error que se termino de enviar el stream
			return stream.SendAndClose(&testpb.SetQuestionResponse{ // Cerramos la sesion y enviamos la respuesta al cliente
				Ok: true,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return err
		}
		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.GetAnswer(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
		}
		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			fmt.Println(err)
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}

	}
}

// STUDENTS
func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	// Iterrar a traves de los mensajes recibidos
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		enrollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())

	if err != nil {
		return err
	}

	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}
		err := stream.Send(student)

		time.Sleep(2 * time.Second) // optional

		if err != nil {
			return err
		}

	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {

	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "test1") // test1 debe ser parametrizado
	if err != nil {
		return err
	}

	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}

			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
			answer, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			log.Println("Answer: ", answer.GetAnswer())
		}
	}
}

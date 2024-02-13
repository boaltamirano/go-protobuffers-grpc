// ESTE ARCHIVO ES EL SERVIDOR DEL MODELO TEST

package server

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/boaltamirano/go-protobuffers-grpc/models"
	"github.com/boaltamirano/go-protobuffers-grpc/repository"
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

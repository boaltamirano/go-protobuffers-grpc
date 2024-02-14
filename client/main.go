package main

import (
	"context"
	"log"
	"time"

	"github.com/boaltamirano/go-protobuffers-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Creamos una configuracion insegura en el segundo parametro del Dial
	// TODO: Investigar como encriptar esa coneccion
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	// Creamos un  cliente
	c := testpb.NewTestServiceClient(cc)

	// DoUnary(c)
	DoClientStreaming(c)

}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "test1",
	}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest: %v", err)
	}

	log.Printf("response from server: %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "question8test1",
			Answer:   "Azul",
			Question: "Color asociado a Goland",
			TestId:   "test1",
		},
		{
			Id:       "question9test1",
			Answer:   "Google",
			Question: "Empresa que desarrollo goland",
			TestId:   "test1",
		},
		{
			Id:       "question10test1",
			Answer:   "Backend",
			Question: "Especialidad de Goland",
			TestId:   "test1",
		},
	}

	stream, err := c.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("error while calling SetQuestions: %v", err)
	}

	for _, question := range questions {
		log.Println("sending question: ", question.Id)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	log.Printf("response from server: %v", msg)
}

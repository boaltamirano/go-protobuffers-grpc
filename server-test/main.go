// ESTE ARCHIVO MAIN ES PARA INICIAL EL SERVIDOR DEL MODELO TEST

package main

import (
	"log"
	"net"

	"github.com/boaltamirano/go-protobuffers-grpc/database"
	"github.com/boaltamirano/go-protobuffers-grpc/server"
	"github.com/boaltamirano/go-protobuffers-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Listener para definir el puerto en el que se va a ejecutar
	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")

	server := server.NewTestServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	// Declaramos s, donde definimos que va a iniciar un nuevo servidor
	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, server) // Definimos un server del model test service

	//Agregamos reflection para proder autocompletar la metadata desde el cleinte postman
	reflection.Register(s)

	// Levantamos el server
	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}

}

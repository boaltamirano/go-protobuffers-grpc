package server

import (
	"context"

	"github.com/boaltamirano/go-protobuffers-grpc/models"
	"github.com/boaltamirano/go-protobuffers-grpc/repository"
	"github.com/boaltamirano/go-protobuffers-grpc/studentpb"
)

type Server struct {
	repo                                        repository.Repository //
	studentpb.UnimplementedStudentServiceServer                       // Este import define que la struct server tentra todas las propiedades de studentpb.UnimplementedStudentServiceServer
}

// Constructor de del struct
func NewStudentServe(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

// Implementacion de GetStudent herredando del strudent.proto
func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	student, err := s.repo.GetStudent(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

// Implementacion de SetStudent herredando del strudent.proto
func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {
	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repo.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &studentpb.SetStudentResponse{
		Id: student.Id,
	}, nil
}

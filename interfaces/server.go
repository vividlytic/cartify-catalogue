package interfaces

import (
	"cartify/catalogue/interfaces/service"

	pb "cartify/catalogue/proto/book"

	"cartify/catalogue/app/usecase/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"cartify/catalogue/domain/repository"
)

type ServerParams struct {
	BookRepository repository.BookRepository
}

func NewServer(params ServerParams) *grpc.Server {
	server := grpc.NewServer()

	bookService := service.NewBookServer(
		book.NewListBooks(params.BookRepository),
		book.NewGetBook(params.BookRepository),
	)

	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	healthpb.RegisterHealthServer(server, hsrv)

	reflection.Register(server)

	pb.RegisterCatalogueServer(server, bookService)

	return server

}

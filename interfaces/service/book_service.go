package service

import (
	usecase "cartify/catalogue/app/usecase/book"
	"context"

	"cartify/catalogue/domain/model"
	pb "cartify/catalogue/proto/book"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BookServer struct {
	listBooks usecase.ListBooks
	getBook   usecase.GetBook
	pb.UnimplementedCatalogueServer
}

func NewBookServer(
	listBooks usecase.ListBooks,
	getBook usecase.GetBook,
) pb.CatalogueServer {
	return &BookServer{
		listBooks: listBooks,
		getBook:   getBook,
	}
}

func (s *BookServer) ListBooks(ctx context.Context, in *emptypb.Empty) (*pb.ListBooksResponse, error) {
	// if md, ok := metadata.FromIncomingContext(ctx); ok {
	// 	fmt.Println(md.Get("authorization"))
	// }
	books, err := s.listBooks(ctx)
	if err != nil {
		return nil, err
	}

	protoBooks := make([]*pb.Book, 0)

	for _, b := range books {
		protoBooks = append(protoBooks, BookToProto(b))
	}

	response := &pb.ListBooksResponse{Books: protoBooks}
	return response, nil
}

func BookToProto(book *model.Book) *pb.Book {
	protoBook := &pb.Book{
		Id:     int32(book.ID),
		Title:  book.Title,
		Author: book.Author,
		Price:  int32(book.Price),
	}

	return protoBook
}

func (s *BookServer) GetBook(ctx context.Context, request *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	params := usecase.GetBooksParams{ID: int(request.Id)}

	book, err := s.getBook(ctx, params)
	if err != nil {
		return nil, err
	}

	protoBook := BookToProto(book)

	return &pb.GetBookResponse{Book: protoBook}, nil
}

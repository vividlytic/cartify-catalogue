package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "cartify/catalogue/proto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// 書籍情報の構造体
type Book struct {
	Id     int
	Title  string
	Author string
	Price  int
}

// DBを使わずに書籍情報を返せるように構造体にデータを保持
var (
	book1 = Book{
		Id:     1,
		Title:  "The Awakening",
		Author: "Kate Chopin",
		Price:  1000,
	}
	book2 = Book{
		Id:     2,
		Title:  "City of Glass",
		Author: "Paul Auster",
		Price:  2000,
	}
	books = []Book{book1, book2}
)

func getBook(i int) Book {
	return books[i-1]
}

// CatalogueServerのサービス実装用インターフェイスの構造体
type server struct {
	pb.UnimplementedCatalogueServer
}

// 自動生成された`catalogue_grpc.pb.go`の`GetBook`インターフェースを実装。
func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	// リクエストで指定されたIDに応じて返す書籍情報を取得
	book := getBook(int(in.Id))

	// レスポンス用のデータを作成
	protoBook := &pb.Book{
		Id:     int32(book.Id),
		Title:  book.Title,
		Author: book.Author,
		Price:  int32(book.Price),
	}

	// レスポンス用のコードを使ってレスポンスを作り返却
	return &pb.GetBookResponse{Book: protoBook}, nil
}

// 自動生成された`catalogue_grpc.pb.go`の`ListBooks`インターフェースを実装。
func (s *server) ListBooks(ctx context.Context, in *emptypb.Empty) (*pb.ListBooksResponse, error) {
	// レスポンス用のデータを作成
	protoBooks := make([]*pb.Book, 0)

	for _, book := range books {
		protoBook := &pb.Book{
			Id:     int32(book.Id),
			Title:  book.Title,
			Author: book.Author,
			Price:  int32(book.Price),
		}
		protoBooks = append(protoBooks, protoBook)
	}

	// レスポンス用のコードを使ってレスポンスを作り返却
	return &pb.ListBooksResponse{Books: protoBooks}, nil
}

var (
	// 起動するポート番号を指定
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// gRPCサーバーの生成
	s := grpc.NewServer()
	// 自動生成された関数に、サーバと実際に処理を行うメソッドを実装したハンドラを設定
	pb.RegisterCatalogueServer(s, &server{})
	// gRPCサーバーにリフレクションサービスを登録
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	// サーバーを起動し、エラー発生時にはエラーメッセージを出力
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

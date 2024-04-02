package main

import (
	"log"
	"net"
	"os"

	"github.com/Nishad4140/SkillSync_ProjectService/db"
	"github.com/Nishad4140/SkillSync_ProjectService/initializer"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/service"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}
	addr := os.Getenv("DB_KEY")
	db, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	userConn, err := grpc.Dial("localhost:4001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while connecting to user service")
	}

	defer func() {
		userConn.Close()
	}()

	userRes := pb.NewUserServiceClient(userConn)

	service.UserClient = userRes

	services := initializer.Initializer(db)
	server := grpc.NewServer()
	pb.RegisterProjectServiceServer(server, services)
	lis, err := net.Listen("tcp", ":4002")
	if err != nil {
		log.Fatalf("failed to run on the port 4002 : %v", err)
	}
	log.Printf("user service listening on the port 4002")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to run on the port 4002 : %v", err)
	}
}

package main

import (
	"context"
	"errors"
	"flag"
	"github.com/CrescentKohana/Zeniire/internal/auth"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"github.com/CrescentKohana/Zeniire/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"net"

	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"google.golang.org/grpc"
)

// server is used to implement RecordsServer.
type server struct {
	pb.UnimplementedRecordsServer
}

var dbAPI db.API

// CreateRecord implements the function for creating and storing Records to the db.
func (*server) CreateRecord(_ context.Context, req *pb.CreateRecordReq) (*pb.CreateRecordResp, error) {
	amount := req.GetAmount()
	datetime := req.GetDatetime()
	id := uuid.New().String()

	data := pb.Record{
		Uuid:     id,
		Amount:   amount,
		Datetime: datetime,
	}

	if err := dbAPI.CreateRecord(&data); err != nil {
		log.Error(err)
		return nil, errors.New("record creation unsuccessful")
	}

	return &pb.CreateRecordResp{
		Record: &data,
	}, nil
}

// ReturnRecord implements the function for reading and returning singular Records from the db.
func (*server) ReturnRecord(_ context.Context, req *pb.ReadRecordReq) (*pb.ReadRecordResp, error) {
	record, err := dbAPI.ReturnRecord(req.GetRecordUuid())

	if err != nil {
		return nil, errors.New("record not found")
	}
	return &pb.ReadRecordResp{
		Record: record,
	}, nil
}

// ReturnRecords implements the function for reading and returning multiple Records from the db.
// Optionally can be given a timerange in the RFC3339 format.
func (*server) ReturnRecords(_ context.Context, req *pb.ReadRecordsReq) (*pb.ReadRecordsResp, error) {
	records, err := dbAPI.ReturnRecords(req.StartDatetime, req.EndDatetime, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return &pb.ReadRecordsResp{
		Records: records,
	}, nil
}

// main launches the server while loading environmentals, initializing the DB and migrations.
func main() {
	// Load environmentals
	config.LoadEnv()

	// Initialize database (PostgreSQL or mock)
	conn, connErr := pgx.Connect(context.Background(), config.Options.DB.Address)
	if connErr != nil {
		// If the database connection was unsuccessful, exit the application with an error.
		log.Fatal(connErr)
	}
	dbAPI = db.API{Db: conn}
	db.EnsureLatestVersion()

	flag.Parse()
	lis, err := net.Listen("tcp", config.Options.GRPC.Host+":"+config.Options.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var s *grpc.Server
	if config.Options.GRPC.TLS {
		tlsCredentials, err := auth.LoadServerTLSCredentials()
		if err != nil {
			log.Fatal("could not load TLS credentials: ", err)
		}
		s = grpc.NewServer(grpc.Creds(tlsCredentials))
	} else {
		s = grpc.NewServer()
	}

	pb.RegisterRecordsServer(s, &server{})
	log.Info("Zeniire server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

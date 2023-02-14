package main

import (
	"context"
	"errors"
	"flag"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"github.com/CrescentKohana/Zeniire/pkg/db"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"

	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"google.golang.org/grpc"
)

// server is used to implement RecordsServer.
type server struct {
	pb.UnimplementedRecordsServer
}

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

	if err := db.CreateRecord(&data); err != nil {
		log.Error(err)
		return nil, errors.New("record creation unsuccessful")
	}

	return &pb.CreateRecordResp{
		Record: &data,
	}, nil
}

// ReturnRecord implements the function for reading and returning singular Records from the db.
func (*server) ReturnRecord(_ context.Context, req *pb.ReadRecordReq) (*pb.ReadRecordResp, error) {
	record, err := db.ReturnRecord(req.GetRecordUuid())

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
	records, err := db.ReturnRecords(req.StartDatetime, req.EndDatetime)
	if err != nil {
		log.Error(err)
		return nil, errors.New("no records found")
	}
	return &pb.ReadRecordsResp{
		Records: records,
	}, nil
}

// main launches the server while loading environmentals, initializing the DB and migrations.
func main() {
	config.LoadEnv()
	db.Initdb()
	db.EnsureLatestVersion()

	flag.Parse()
	lis, err := net.Listen("tcp", config.Options.GRPC.Host+":"+config.Options.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRecordsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

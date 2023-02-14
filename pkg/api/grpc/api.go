package grpcAPI

import (
	"context"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"github.com/CrescentKohana/Zeniire/pkg/utility"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var conn *grpc.ClientConn
var client pb.RecordsClient

func ReturnRecord(id string) (*pb.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.ReturnRecord(ctx, &pb.ReadRecordReq{RecordUuid: id})
	return r.GetRecord(), err
}

func ReturnRecords(startDatetime string, endDatetime string) ([]*pb.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.ReturnRecords(ctx, &pb.ReadRecordsReq{
		EndDatetime:   utility.StringToTimestamp(startDatetime),
		StartDatetime: utility.StringToTimestamp(endDatetime),
	})

	return r.GetRecords(), err
}

func CreateRecord(amount int64, datetime *timestamppb.Timestamp) (*pb.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.CreateRecord(ctx, &pb.CreateRecordReq{
		Amount:   amount,
		Datetime: datetime,
	})

	return r.GetRecord(), err
}

func InitGRPCClient() {
	var err error
	conn, err = grpc.Dial(
		config.Options.GRPC.Host+":"+config.Options.GRPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}

	// defer conn.Close()
	client = pb.NewRecordsClient(conn)
}

package grpcAPI

import (
	"context"
	"github.com/CrescentKohana/Zeniire/internal/auth"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"github.com/CrescentKohana/Zeniire/pkg/utility"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func ReturnRecords(startDatetime string, endDatetime string, limit int64, offset int64) ([]*pb.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.ReturnRecords(ctx, &pb.ReadRecordsReq{
		StartDatetime: utility.StringToTimestamp(startDatetime),
		EndDatetime:   utility.StringToTimestamp(endDatetime),
		Limit:         limit,
		Offset:        offset,
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
	var connErr error
	var tlsCredentials credentials.TransportCredentials
	if config.Options.GRPC.TLS {
		var err error
		tlsCredentials, err = auth.LoadClientTLSCredentials()
		if err != nil {
			log.Fatal("could not load TLS credentials: ", err)
		}
	} else {
		tlsCredentials = insecure.NewCredentials()
	}

	conn, connErr = grpc.Dial(
		config.Options.GRPC.Host+":"+config.Options.GRPC.Port,
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if connErr != nil {
		log.Fatal("could not dial the gRPC server: ", connErr)
	}

	// defer conn.Close()
	client = pb.NewRecordsClient(conn)
}

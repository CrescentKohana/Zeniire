package db

import (
	"context"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ReturnRecords(startdateime *timestamppb.Timestamp, enddatetime *timestamppb.Timestamp) ([]*pb.Record, error) {
	var rows pgx.Rows
	if startdateime == nil && enddatetime == nil {
		rows, _ = conn.Query(context.Background(), "SELECT * FROM records")
	} else {
		if startdateime == nil {
			startdateime = timestamppb.New(time.Unix(0, 0))
		}
		if enddatetime == nil {
			startdateime = timestamppb.Now()
		}
		log.Info(startdateime.AsTime())
		// BETWEEN should not be used as it will include results where the timestamp is exactly 2023-01-01 00:00:00.000000,
		// but not timestamps later in that same day.
		rows, _ = conn.Query(
			context.Background(),
			"SELECT * FROM records WHERE datetime >= $1 AND datetime < $2",
			startdateime.AsTime(),
			enddatetime.AsTime(),
		)
	}

	var records []*pb.Record
	for rows.Next() {
		var id string
		var datetime *time.Time
		var amount int64

		if err := rows.Scan(&id, &amount, &datetime); err != nil {
			return nil, err
		}

		records = append(records, &pb.Record{
			Uuid:     id,
			Amount:   amount,
			Datetime: timestamppb.New(*datetime),
		})
	}

	return records, rows.Err()
}

func ReturnRecord(id string) (*pb.Record, error) {
	var uuid string
	var amount int64
	var datetime time.Time

	err := conn.
		QueryRow(context.Background(), "select uuid, amount, datetime from records where uuid=$1", id).
		Scan(&uuid, &amount, &datetime)

	return &pb.Record{
		Uuid:     uuid,
		Amount:   amount,
		Datetime: timestamppb.New(datetime),
	}, err
}

func CreateRecord(data *pb.Record) error {
	_, err := conn.Exec(
		context.Background(),
		"insert into records(uuid, amount, datetime) values($1, $2, $3)", data.Uuid, data.Amount, data.Datetime.AsTime(),
	)
	return err
}

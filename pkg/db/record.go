package db

import (
	"context"
	"github.com/CrescentKohana/Zeniire/pkg/utility"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"time"
)

func (c *API) ReturnRecords(startDatetime *timestamppb.Timestamp, endDatetime *timestamppb.Timestamp, limit int64, offset int64) ([]*pb.Record, error) {

	limit = utility.Clamp(limit, 1, 1000)
	offset = utility.Clamp(offset, 0, math.MaxInt64)

	var rows pgx.Rows
	if startDatetime == nil && endDatetime == nil {
		var err error
		rows, err = c.Db.Query(context.Background(), "SELECT * FROM records LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		if startDatetime == nil {
			startDatetime = timestamppb.New(time.Unix(0, 0))
		}
		if endDatetime == nil {
			endDatetime = timestamppb.Now()
		}

		// BETWEEN should not be used as it will include results where the timestamp is exactly 2023-01-01 00:00:00.000000,
		// but not timestamps later in that same day.
		var err error
		rows, err = c.Db.Query(
			context.Background(),
			"SELECT * FROM records WHERE datetime >= $1 AND datetime < $2  LIMIT $3 OFFSET $4",
			startDatetime.AsTime(),
			endDatetime.AsTime(),
			limit,
			offset,
		)

		if err != nil {
			log.Error(err)
			return nil, err
		}
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

func (c *API) ReturnRecord(id string) (*pb.Record, error) {
	var uuid string
	var amount int64
	var datetime time.Time

	err := c.Db.
		QueryRow(context.Background(), "SELECT uuid, amount, datetime FROM records WHERE uuid=$1", id).
		Scan(&uuid, &amount, &datetime)

	return &pb.Record{
		Uuid:     uuid,
		Amount:   amount,
		Datetime: timestamppb.New(datetime),
	}, err
}

func (c *API) CreateRecord(data *pb.Record) error {
	_, err := c.Db.Exec(
		context.Background(),
		"INSERT INTO records(uuid, amount, datetime) VALUES($1, $2, $3)", data.Uuid, data.Amount, data.Datetime.AsTime(),
	)
	return err
}

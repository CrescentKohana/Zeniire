package restAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CrescentKohana/Zeniire/internal/config"
	grpcAPI "github.com/CrescentKohana/Zeniire/pkg/api/grpc"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ReturnRecordsForm struct {
	StartDatetime string
	EndDatetime   string
}

// resultToJSON returns the given struct as JSON format for HTTP responses.
func resultToJSON(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "code": %d, "message": "error while encoding JSON" }`, http.StatusInternalServerError)
	}
}

func RecordCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "recordID")
		record, err := grpcAPI.ReturnRecord(id)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "record", record)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	record, ok := ctx.Value("record").(*pb.Record)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.Write([]byte(fmt.Sprintf("record: %s", record.Datetime)))
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	formData := &ReturnRecordsForm{}
	if err := json.NewDecoder(r.Body).Decode(formData); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	records, err := grpcAPI.ReturnRecords(formData.StartDatetime, formData.EndDatetime)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.Write([]byte(fmt.Sprintf("records: %s", records)))
}

func createRecord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recordData, ok := ctx.Value("record").(*pb.Record)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	record, err := grpcAPI.CreateRecord(recordData.Amount, recordData.Datetime)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	w.Write([]byte(fmt.Sprintf("created: %s", record.Datetime)))
}

// InitRESTClient initializes the REST HTTP client.
func InitRESTClient() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resultToJSON(w,
			struct {
				Name       string
				Version    string
				APIVersion int
				Git        string
			}{
				Name:       config.AppName,
				Version:    config.AppVersion,
				APIVersion: config.APIVersion,
				Git:        config.GitLink,
			})
	})

	// RESTy routes for the Records resource
	r.Route("/records", func(r chi.Router) {
		r.Get("/", getRecords)    // GET /records
		r.Post("/", createRecord) // POST /records

		// Subrouters
		r.Route("/{recordID}", func(r chi.Router) {
			r.Use(RecordCtx)
			r.Get("/", getRecord) // GET /records/f982e6d6-a28d-406d-9378-dce5972ae6d5
		})
	})

	if err := http.ListenAndServe(config.Options.RESTAddr, r); err != nil {
		log.Error(err)
		return
	}
}

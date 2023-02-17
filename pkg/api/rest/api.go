package restAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CrescentKohana/Zeniire/internal/config"
	grpcAPI "github.com/CrescentKohana/Zeniire/pkg/api/grpc"
	"github.com/CrescentKohana/Zeniire/pkg/utility"
	pb "github.com/CrescentKohana/Zeniire/proto/gen/go/zeniire"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type NewRecordForm struct {
	Amount   int64
	Datetime string
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

	resultToJSON(w, record)
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	startDatetime := r.URL.Query().Get("startDatetime")
	endDatetime := r.URL.Query().Get("endDatetime")
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 1000
	}
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	records, err := grpcAPI.ReturnRecords(startDatetime, endDatetime, limit, offset)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	resultToJSON(w, records)
}

func createRecord(w http.ResponseWriter, r *http.Request) {
	formData := &NewRecordForm{}
	if err := json.NewDecoder(r.Body).Decode(formData); err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if formData.Amount < 0 {
		http.Error(w, "negative amounts not allowed", 400)
		return
	}

	parsedDatetime := utility.StringToTimestamp(formData.Datetime)
	if parsedDatetime == nil {
		http.Error(w, "invalid timestamp", 400)
		return
	}

	record, err := grpcAPI.CreateRecord(formData.Amount, parsedDatetime)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	resultToJSON(w, record)
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

	log.Info("Starting REST API on: ", config.Options.RESTAddr)
	if err := http.ListenAndServe(config.Options.RESTAddr, r); err != nil {
		log.Error(err)
		return
	}
}

package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Bronku/iroon/models"
	"github.com/Bronku/iroon/store"
	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type config struct {
	Database struct {
		File string
	}

	Server struct {
		Addr              string
		ReadHeaderTimeout time.Duration
		RequestTimeout    time.Duration
	}
}

//go:embed config.default.toml
var defaultConfig string

type contextKey string

const dbContextKey = contextKey("db")
const orderContextKey = contextKey("order")

func main() {
	// load config
	var conf config
	_, err := toml.Decode(defaultConfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Println("using default config")
	}

	// load database
	db, err := store.LoadStore(conf.Database.File)
	if err != nil {
		log.Panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(conf.Server.RequestTimeout))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.WithValue(dbContextKey, db))
		r.Route("/{orderID}", func(r chi.Router) {
			r.Use(OrderCtx) // Load the *Article on the request context
			r.Get("/", getOrder)
			// r.Get("/", GetArticle)       // GET /articles/123
			// r.Put("/", UpdateArticle)    // PUT /articles/123
			// r.Delete("/", DeleteArticle) // DELETE /articles/123
		})

	})

	// start server
	server := &http.Server{
		Handler:           r,
		Addr:              conf.Server.Addr,
		ReadHeaderTimeout: conf.Server.ReadHeaderTimeout,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// can't get item from context, please include it before registering route
var errNoContextValue = errors.New("no value from context")

func OrderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := r.Context().Value(dbContextKey).(*gorm.DB)
		if db == nil {
			panic(errNoContextValue)
		}
		var order *models.Order
		var err error

		orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
		if err != nil {
			w.Write([]byte("url not found"))
			return
		}

		err = db.Preload("OrderItems.Product").Preload(clause.Associations).Find(&order, orderID).Error
		if err != nil {
			w.Write([]byte("item not found"))
			return
		}

		ctx := context.WithValue(r.Context(), orderContextKey, order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	var db *gorm.DB
	var catalogue []models.Product
	var order *models.Order
	var err error

	db = r.Context().Value(dbContextKey).(*gorm.DB)
	if db == nil {
		panic(errNoContextValue)
	}

	err = db.Find(&catalogue).Error
	if err != nil {
		w.Write([]byte("no catalogue"))
		return
	}

	order = r.Context().Value(orderContextKey).(*models.Order)
	if order == nil {
		panic(errNoContextValue)
	}

	fmt.Fprint(w, "order: ", order, "\ncatalogue: ", catalogue)
}

package main

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Bronku/iroon/models"
	"github.com/Bronku/iroon/store"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:embed static/*
var static embed.FS

type contextKey string

const dbContextKey = contextKey("db")
const orderContextKey = contextKey("order")

func main() {
	loadConfig()
	loadTemplates()

	// load database
	db, err := store.LoadStore(config.Database.File)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(config.Server.RequestTimeout))

	r.Get("/", http.RedirectHandler("/orders", http.StatusSeeOther).ServeHTTP)

	// Serve static files from the "./assets" directory when accessing "/assets/*"
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(static)))
	r.Handle("/static/*", fileServer)

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.WithValue(dbContextKey, db))
		r.Route("/{orderID}", func(r chi.Router) {
			r.Use(OrderCtx)
			r.Get("/", getOrder)
		})
	})

	// start server
	server := &http.Server{
		Handler:           r,
		Addr:              config.Server.Addr,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
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
			log.Fatal(errNoContextValue)
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
		log.Fatal(errNoContextValue)
	}

	err = db.Find(&catalogue).Error
	if err != nil {
		w.Write([]byte("no catalogue"))
		return
	}

	order = r.Context().Value(orderContextKey).(*models.Order)
	if order == nil {
		log.Fatal(errNoContextValue)
	}

	templates["order.gohtml"].ExecuteTemplate(w, "layout", struct {
		Order     models.Order
		Catalogue []models.Product
	}{*order, catalogue})
}

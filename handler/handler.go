package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type service interface {
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenURL(ctx context.Context, url string) (string, error)
}

type Handler struct {
	serv service
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.GetURL)
	mux.HandleFunc("/api/shorten", h.ShortenURL)
	return mux
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if r.Method == http.MethodGet {
		log.Println(r.URL)
		shortURL := r.URL.String()

		// trimming / from the start of the url
		shortURL = shortURL[1:]
		log.Println("ShortURL: ", shortURL)

		url, err := h.serv.GetURL(ctx, shortURL)
		if err != nil {
			fmt.Println("Server Error: ", err.Error())
			serverError(w, err.Error())
			return
		}

		log.Println("URL: ", url);
		http.Redirect(w, r, url, http.StatusSeeOther)
	}  else {
		clientError(w, "Method not supported.")
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			clientError(w, "Content-Type header should be set text/plain; charset=utf-8")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			serverError(w, err.Error())
		}

		url := string(body)

		shortURL, err := h.serv.ShortenURL(ctx, url)
		if err != nil {
			serverError(w, err.Error())
		}

		w.Write([]byte(shortURL))
	} else {
		clientError(w, "Method not supported.")
	}
}

func serverError(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func NewHandler(s service) *Handler {

	return &Handler{
		serv: s,
	}
}
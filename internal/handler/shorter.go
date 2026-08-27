package handler

import (
	"io"
	"net/http"
)

type ShorterService interface {
	DoShortUrl(url string) string
	GetUrl(shortUrl string) (string, error)
}

type Server struct {
	service ShorterService
}

func NewServer(s ShorterService) *Server {
	return &Server{service: s}
}

func (s *Server) DoShortUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := string(body)
	if url == "" {
		http.Error(w, "пустой URL", http.StatusBadRequest)
		return
	}

	shortUrl := s.service.DoShortUrl(url)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func (s *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	urlId := r.PathValue("id")
	url, err := s.service.GetUrl(urlId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

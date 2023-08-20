package url

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/neglarken/url_shortener/internal/lib/api/random"
	resp "github.com/neglarken/url_shortener/internal/lib/api/response"
	"github.com/neglarken/url_shortener/internal/storage"
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	DeleteURL(alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

func New(log *log.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Print("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Println("failed to decode request body", err.Error())

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Print("request body decoded", req)

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Println("invalid request", err.Error())

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Println("url already exists", req.URL)

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Println("failed to add url", err.Error())

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Println("url added", id)

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

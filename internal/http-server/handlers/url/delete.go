package url

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/neglarken/url_shortener/internal/lib/api/response"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

type RequestDelete struct {
	Alias string `json:"alias,omitempty"`
}

type ResponseDelete struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func Delete(log *log.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Delete"

		var req RequestDelete

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

		if err = urlSaver.DeleteURL(req.Alias); err != nil {
			log.Println("invalid request", err.Error())
			return
		}

		log.Println("alias deleted", req.Alias)

		responseOK(w, r, req.Alias)
	}
}

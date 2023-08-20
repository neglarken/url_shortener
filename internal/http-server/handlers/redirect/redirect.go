package redirect

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "github.com/neglarken/url_shortener/internal/lib/api/response"
	"github.com/neglarken/url_shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *log.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Println("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Println("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Println("failed to get url", err.Error())

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Println("got url", resURL)

		// redirect to found url
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}

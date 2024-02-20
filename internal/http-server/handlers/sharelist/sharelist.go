package sharelist

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"wish_list/internal/entity"
	resp "wish_list/internal/lib/api/response"
)

type ShareList interface {
	GetList(alias string) ([]entity.GiftList, error)
}

func GetList(log *slog.Logger, shareList ShareList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sharelist.GetList"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias parameter is required"))
			return
		}

		list, err := shareList.GetList(alias)
		if err != nil {
			log.Error("failed to get list: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("unable to retrieve the wishlist"))
			return
		}

		if list == nil {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("wishlist not found"))
			return
		}

		render.JSON(w, r, list)
	}
}

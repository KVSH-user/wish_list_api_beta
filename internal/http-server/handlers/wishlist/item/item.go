package item

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"wish_list/internal/entity"
	"wish_list/internal/http-server/handlers/auth/uidextractor"
	resp "wish_list/internal/lib/api/response"
)

type Request struct {
	WishListId int    `json:"wish_list_id"`
	GiftName   string `json:"gift_name"`
	Url        string `json:"url"`
}

type Response struct {
	GiftId int `json:"gift_id"`
}

type Item interface {
	CreateItem(wishlistId int, giftName, url string) (int, error)
	GetByWishId(wishListId int) ([]entity.GiftList, error)
	DelItemById(itemId int) error
}

func Create(log *slog.Logger, item Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.item.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		uid, token, err := uidextractor.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, _ = token, uid

		var req Request

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		giftId, err := item.CreateItem(req.WishListId, req.GiftName, req.Url)
		if err != nil {
			log.Error("failed to add item: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("item added")

		render.JSON(w, r, Response{
			GiftId: giftId,
		})

	}
}

func GetByWishId(log *slog.Logger, item Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.item.GetByWishId"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		uid, token, err := uidextractor.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, _ = token, uid

		wishListId := chi.URLParam(r, "wishlistId")
		if wishListId == "" {
			log.Info("wishListId is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("wishListId parameter is required"))
			return
		}

		wishListIdInt, err := strconv.Atoi(wishListId)
		if err != nil {
			log.Error("Invalid wishListId format: ", err)
			http.Error(w, "Invalid wishListId", http.StatusBadRequest)
			return
		}

		lists, err := item.GetByWishId(wishListIdInt)
		if err != nil {
			log.Error("failed to get items: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("items got")

		render.JSON(w, r, lists)

	}
}

type RequestForDel struct {
	ItemId int `json:"item_id"`
}

func Delete(log *slog.Logger, item Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.item.Delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		uid, token, err := uidextractor.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_, _ = token, uid

		var req RequestForDel

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = item.DelItemById(req.ItemId)
		if err != nil {
			log.Error("failed to delete item: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("item deleted")

		render.JSON(w, r, http.StatusOK)

	}
}

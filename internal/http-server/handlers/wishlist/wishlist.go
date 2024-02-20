package wishlist

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"wish_list/internal/entity"
	"wish_list/internal/http-server/handlers/auth/uidextractor"
	alias2 "wish_list/internal/lib/alias"
	resp "wish_list/internal/lib/api/response"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	WishListId int    `json:"wish_list_id"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`
}

type Wishlist interface {
	CreateList(name, alias string, uid int) (int, error)
	GetAllLists(uid int) ([]entity.WishList, error)
	WishListDel(wishListId int) error
}

func Create(log *slog.Logger, wishlist Wishlist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wishlist.Create"

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

		_ = token

		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			log.Error("Invalid UID format: ", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

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

		alias := alias2.NewRandomString(uidInt)

		wishListId, err := wishlist.CreateList(req.Name, alias, uidInt)
		if err != nil {
			log.Error("failed to create wishlist: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("wishlist created")

		render.JSON(w, r, Response{
			WishListId: wishListId,
			Name:       req.Name,
			Alias:      alias,
		})

	}
}

func GetAllLists(log *slog.Logger, wishlist Wishlist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wishlist.GetAllLists"

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

		_ = token

		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			log.Error("Invalid UID format: ", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		lists, err := wishlist.GetAllLists(uidInt)
		if err != nil {
			log.Error("failed to get wishlists: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("wishlists got")

		render.JSON(w, r, lists)

	}
}

type RequestForDel struct {
	WishListId int `json:"wish_list_id"`
}

func Delete(log *slog.Logger, wishlist Wishlist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wishlist.Create"

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

		err = wishlist.WishListDel(req.WishListId)
		if err != nil {
			log.Error("failed to delete wishlist: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("wishlist deleted")

		render.JSON(w, r, http.StatusOK)

	}
}

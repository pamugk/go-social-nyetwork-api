package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/pamugk/social-nyetwork-server/internal/app"
)

func StartServer(port *int) {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/"))

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(middleware.Timeout(time.Second * 60))

		r.Route("/user", func(r chi.Router) {
			r.With(pagination).Get("/", func(w http.ResponseWriter, r *http.Request) {
				loginPart := r.URL.Query().Get("loginPart")
				page := r.Context().Value("page").(int32)
				limit := r.Context().Value("limit").(int32)
				pageContents, total, err := app.SearchUsers(loginPart, page, limit)
				if err != nil {
					w.WriteHeader(400)
					render.Render(w, r, &errorResponse{err.Error()})
					return
				}

				response := &searchUsersResponse{Page: pageContents}
				response.PageNumber = page
				response.PageSize = limit
				response.Total = total
				render.Render(w, r, response)

			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				request := &createUserRequest{}
				if err := render.Bind(r, request); err != nil {
					w.WriteHeader(400)
					render.Render(w, r, &errorResponse{err.Error()})
					return
				}
				data, err := convertFromUserData(&request.UserData)
				if err == nil {
					_, err = app.CreateUser(data, request.Password)
				}
				if err != nil {
					w.WriteHeader(400)
					render.Render(w, r, &errorResponse{err.Error()})
					return
				}

				w.WriteHeader(201)
			})

			r.Route("/{id:^-?\\d+$}", func(r chi.Router) {
				r.Use(idContext)
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					userId := r.Context().Value("id").(int64)
					user, err := app.GetUser(userId)
					if err == nil {
						render.Render(w, r, &getUserResponse{*convertToUser(user)})
					} else if err.Error() == "Not found" {
						w.WriteHeader(404)
					} else {
						w.WriteHeader(400)
						render.Render(w, r, &errorResponse{err.Error()})
					}
				})
				r.Put("/", func(w http.ResponseWriter, r *http.Request) {
					userId := r.Context().Value("id").(int64)
					request := &updateUserRequest{}
					if err := render.Bind(r, request); err != nil {
						w.WriteHeader(400)
						render.Render(w, r, &errorResponse{err.Error()})
						return
					}
					data, err := convertFromUserData(&request.UserData)
					if err != nil {
						err = app.UpdateUser(userId, data)
					}
					if err != nil {
						if err.Error() == "Not found" {
							w.WriteHeader(404)
						} else {
							w.WriteHeader(400)
							render.Render(w, r, &errorResponse{err.Error()})
						}
					}
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					userId := r.Context().Value("id").(int64)
					if err := app.DeleteUser(userId); err != nil {
						if err.Error() == "Not found" {
							w.WriteHeader(404)
						} else {
							w.WriteHeader(400)
							render.Render(w, r, &errorResponse{err.Error()})
						}
					}
				})
				r.Put("/password", func(w http.ResponseWriter, r *http.Request) {
					userId := r.Context().Value("id").(int64)
					request := &changePasswordRequest{}
					if err := render.Bind(r, request); err != nil {
						w.WriteHeader(400)
						render.Render(w, r, &errorResponse{err.Error()})
						return
					}
					if err := app.ChangePassword(userId, request.NewPassword); err != nil {
						if err.Error() == "Not found" {
							w.WriteHeader(404)
						} else {
							w.WriteHeader(400)
							render.Render(w, r, &errorResponse{err.Error()})
						}
					}
				})
			})
		})
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), r); err != nil {
		log.Fatalf("Failed to serve REST API: %v", err)
	}
}

func idContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			w.WriteHeader(400)
			render.Render(w, r, &errorResponse{err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var page, limit int64
		limit = 100
		if r.URL.Query().Has("page") {
			page, _ = strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
		}
		if r.URL.Query().Has("limit") {
			limit, _ = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
		}

		ctx := context.WithValue(r.Context(), "page", int32(page))
		ctx = context.WithValue(ctx, "limit", int32(limit))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cur *createUserRequest) Bind(r *http.Request) error {
	return nil
}

func (uur *updateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (cpr *changePasswordRequest) Bind(r *http.Request) error {
	return nil
}

func (e *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	oplog := httplog.LogEntry(r.Context())
	oplog.Error().Msg(e.Err)
	return nil
}

func (u *getUserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *searchUsersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

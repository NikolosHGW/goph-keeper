package router

import (
	"github.com/NikolosHGW/goph-keeper/internal/handler"
	"github.com/go-chi/chi"
)

func NewRouter(handlers *handler.Handlers) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(middlewares.Logger.WithLogging)
	// r.Use(middlewares.Gzip.WithGzip)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.AuthHandler.RegisterUser)
		// r.Post("/login", handlers.UserHandler.LoginUser)

		// r.With(middlewares.Auth.WithAuth).Post("/orders", handlers.OrderHandler.UploadOrder)
		// r.With(middlewares.Auth.WithAuth).Get("/orders", handlers.OrderHandler.GetOrders)
		// r.With(middlewares.Auth.WithAuth).Get("/balance", handlers.BalanceHandler.GetBalance)
		// r.With(middlewares.Auth.WithAuth).Post("/balance/withdraw", handlers.WithdrawalHandler.Withdraw)
		// r.With(middlewares.Auth.WithAuth).Get("/withdrawals", handlers.WithdrawalHandler.GetWithdrawals)
	})

	return r
}

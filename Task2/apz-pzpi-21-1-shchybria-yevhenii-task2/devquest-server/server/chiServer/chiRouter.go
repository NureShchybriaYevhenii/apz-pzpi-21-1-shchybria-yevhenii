package chiServer

import (
	"devquest-server/server/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)
func getRoutes() http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.EnableCORS)
	
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	mux.Route("/auth", func(r chi.Router) {
		authSettings := GetChiServer().AuthSettings

		r.Post("/login", authHttpHandler.Login(authSettings))
		r.Post("/register", authHttpHandler.Register(authSettings))
		r.Delete("/logout", authHttpHandler.Logout(authSettings))
	})
	
	mux.Route("/companies", func(r chi.Router) {
		r.Get("/", companyHttpHandler.GetAllCompanies)
		r.Get("/{id}", companyHttpHandler.GetCompanyByID)
		r.Post("/", companyHttpHandler.AddCompany)
		r.Put("/{id}", companyHttpHandler.UpdateCompany)
		r.Delete("/{id}", companyHttpHandler.DeleteCompany)
	})

	return mux
}
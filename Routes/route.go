package Routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
	"github.com/om00/golang-ecommerce/Middleware"
)

type Api struct {
	ctrl   *Controllers.Controller
	router *mux.Router
}

func (api *Api) UserRoutes() {

	api.router.HandleFunc("/users/sign-up", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.Signup(w, r)
	}).Methods("POST")
	api.router.HandleFunc("/users/login", api.ctrl.Login).Methods("POST")
	// api.router.HandleFunc("/admin/addproduct", func(w http.ResponseWriter, r *http.Request) {
	// 	api.ctrl.ProductViewAdmin(w, r)
	// }).Methods("POST")
	// api.router.HandleFunc("/users/prdouctview", func(w http.ResponseWriter, r *http.Request) {
	// 	api.ctrl.SearchProduct(w, r)
	// }).Methods("GET")
	api.router.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.SearchProductByQuery(w, r)
	}).Methods("GET")
	api.router.HandleFunc("/users/refreshToken", api.ctrl.RefreshToken).Methods("GET")
}

func (api *Api) RoutesWithMiddleWare() {

	authRouter := api.router.PathPrefix("/").Subrouter()
	authRouter.Use(Middleware.Authentication)
	authRouter.HandleFunc("/addAddress", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.AddAddress(w, r)
	}).Methods("POST")
	authRouter.HandleFunc("/DeleteAddress", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.DeleteAddress(w, r)
	}).Methods("GET")
	authRouter.HandleFunc("/EditAddress", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.EditAddress(w, r)
	}).Methods("PUT")
	authRouter.HandleFunc("/ListCart", func(w http.ResponseWriter, r *http.Request) {
		api.ctrl.GetItemFromCart(w, r)
	}).Methods("GET")
	authRouter.HandleFunc("/addtocart", api.ctrl.AddToCart).Methods("GET")
	authRouter.HandleFunc("/addtocart", api.ctrl.AddToCart).Methods("GET")
	authRouter.HandleFunc("/removeitem", api.ctrl.RemoveItem).Methods("GET")
	authRouter.HandleFunc("/cartcheckout", api.ctrl.BuyFromCart).Methods("GET")
	authRouter.HandleFunc("/instantbuy", api.ctrl.InstantBuy).Methods("GET")
}

func NewApi(ctrl *Controllers.Controller, router *mux.Router) *Api {
	return &Api{ctrl: ctrl, router: router}
}

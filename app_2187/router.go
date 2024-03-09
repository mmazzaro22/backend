package main

import (
	"github.com/gorilla/mux"
)

// configRouter sets up http server router.
func configRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/action_20071", Action20071Endpoint).Methods("GET")

	router.HandleFunc("/v1/action_21102", Action21102Endpoint).Methods("GET")

	router.HandleFunc("/v1/action_21152", Action21152Endpoint).Methods("POST")

	router.HandleFunc("/v1/action_21261", Action21261Endpoint).Methods("GET")

	router.HandleFunc("/v1/action_21335", Action21335Endpoint).Methods("PUT")

	router.HandleFunc("/v1/action_21369", Action21369Endpoint).Methods("DELETE")

	router.HandleFunc("/v1/action_21581", Action21581Endpoint).Methods("GET")

	router.HandleFunc("/v1/amenities", GetAmenitiesEndpoint).Methods("GET")

	router.HandleFunc("/v1/bookings", CreateBookingsEndpoint).Methods("POST")

	router.HandleFunc("/v1/bookings/listings/{creator_id}/", GetBookingsEndpoint).Methods("GET")

	router.HandleFunc("/v1/bookings/{id}/", GetOneBookingsEndpoint).Methods("GET")

	router.HandleFunc("/v1/bookings/{id}/", DeleteBookingsEndpoint).Methods("DELETE")

	router.HandleFunc("/v1/conversations/{receiver_id}/", Action20241Endpoint).Methods("GET")

	router.HandleFunc("/v1/conversations/{receiver_id}/{sender_id}/", Action21415Endpoint).Methods("GET")

	router.HandleFunc("/v1/create_conversation/{receiver_id}/{sender_id}/{message}", Action20240Endpoint).Methods("POST")

	router.HandleFunc("/v1/help_posts", GetHelpPostsEndpoint).Methods("GET")

	router.HandleFunc("/v1/help_posts/{id}/", GetOneHelpPostsEndpoint).Methods("GET")

	router.HandleFunc("/v1/listings", Action19769Endpoint).Methods("GET")

	router.HandleFunc("/v1/listings/partner_site_id/{partner_site_id}/", Action20255Endpoint).Methods("GET")

	router.HandleFunc("/v1/listings/partner_site_id/{state}/", Action21608Endpoint).Methods("GET")

	router.HandleFunc("/v1/listings/properties/{property_id}", Action20042Endpoint).Methods("GET")

	router.HandleFunc("/v1/listings/{id}/", GetOneListingsEndpoint).Methods("GET")

	router.HandleFunc("/v1/listings/{id}/", DeleteListingsEndpoint).Methods("DELETE")

	router.HandleFunc("/v1/listings/{property_id}/", RequireLogin(Action18697Endpoint)).Methods("POST")

	router.HandleFunc("/v1/logout", LogoutEndpoint).Methods("GET")

	router.HandleFunc("/v1/me", RequireLogin(GetMeEndpoint)).Methods("GET")

	router.HandleFunc("/v1/order_types", GetOrderTypesEndpoint).Methods("GET")

	router.HandleFunc("/v1/partner_sites", GetPartnerSitesEndpoint).Methods("GET")

	router.HandleFunc("/v1/profile_picture", Action21367Endpoint).Methods("GET")

	router.HandleFunc("/v1/properties", RequireLogin(GetPropertiesEndpoint)).Methods("GET")

	router.HandleFunc("/v1/properties", RequireLogin(CreatePropertiesEndpoint)).Methods("POST")

	router.HandleFunc("/v1/properties/all", Action19763Endpoint).Methods("GET")

	router.HandleFunc("/v1/properties/listings/{property_id}/", GetListingsEndpoint).Methods("GET")

	router.HandleFunc("/v1/properties/property_images/{creator_id}", Action21475Endpoint).Methods("GET")

	router.HandleFunc("/v1/properties/{creator_id}", RequireLogin(Action21246Endpoint)).Methods("GET")

	router.HandleFunc("/v1/properties/{id}/", GetOnePropertiesEndpoint).Methods("GET")

	router.HandleFunc("/v1/properties/{id}/", UpdatePropertiesEndpoint).Methods("PUT")

	router.HandleFunc("/v1/properties/{id}/", DeletePropertiesEndpoint).Methods("DELETE")

	router.HandleFunc("/v1/property_images/", Action21272Endpoint).Methods("GET")

	router.HandleFunc("/v1/property_images/{property_id}/", GetOnePropertyImagesEndpoint).Methods("GET")

	router.HandleFunc("/v1/purchase/listing/{id}/", Action21258Endpoint).Methods("GET")

	router.HandleFunc("/v1/request_password_reset", RequestPasswordResetEndpoint).Methods("GET")

	router.HandleFunc("/v1/reset_password", ResetPasswordEndpoint).Methods("GET")

	router.HandleFunc("/v1/signup", SignUpEndpoint).Methods("POST")

	router.HandleFunc("/v1/update_bookings/{id}/", Action20145Endpoint).Methods("PUT")

	router.HandleFunc("/v1/update_bookings/{id}/{reject}/", Action20159Endpoint).Methods("PUT")

	router.HandleFunc("/v1/update_listings/{id}/", UpdateListingsEndpoint).Methods("PUT")

	router.HandleFunc("/v1/update_user", RequireLogin(UpdateUserEndpoint)).Methods("PUT")

	router.HandleFunc("/v1/us_states", GetUsStatesEndpoint).Methods("GET")

	router.HandleFunc("/v1/users", GetUsersEndpoint).Methods("GET")

	router.HandleFunc("/v1/users/{id}/", GetUserEndpoint).Methods("GET")

	return router
}

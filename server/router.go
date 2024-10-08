package server

import (
	_ "github.com/chaaaeeee/sireng/docs"
	trackerHandler "github.com/chaaaeeee/sireng/internal/tracker/handler"
	userHandler "github.com/chaaaeeee/sireng/internal/user/auth/handler"
	userProfileHandler "github.com/chaaaeeee/sireng/internal/user/profile/handler"
	ws "github.com/chaaaeeee/sireng/internal/ws"
	"github.com/chaaaeeee/sireng/middleware"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"net/http"
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	httpSwagger.WrapHandler(w, r)
}

// i js be wrapping shit
func AuthUser(h func(w http.ResponseWriter, r *http.Request), m middleware.Middleware) http.Handler {
	return m.IsUser(m.Authenticate(http.HandlerFunc(h)))
}

func initializeRoutes(mux *http.ServeMux, userHandler userHandler.UserHandler, userProfileHandler userProfileHandler.UserProfileHandler, trackerHandler trackerHandler.TrackerHandler, wsHandler ws.WsHandler, middleware middleware.Middleware) *http.ServeMux {
	// initialize handlers
	mux.HandleFunc("POST /signup", userHandler.SignUp)
	mux.HandleFunc("POST /login", userHandler.Login)
	// mux.HandleFunc("GET /logout", userHandler.Logout)

	mux.Handle("POST /createStudySession", AuthUser(trackerHandler.CreateStudySession, middleware))
	mux.Handle("GET /endStudySession/{userId}", AuthUser(trackerHandler.EndStudySession, middleware))
	mux.Handle("GET /getStudySessions", AuthUser(trackerHandler.GetStudySessions, middleware))
	mux.Handle("GET /getStudySessions/{userId}", AuthUser(trackerHandler.GetStudySessionsByUserId, middleware))

	// User Profile
	mux.Handle("PATCH /updateUsername", AuthUser(userProfileHandler.UpdateUsername, middleware))
	mux.Handle("PATCH /updateProfilePhoto", AuthUser(userProfileHandler.UpdateProfilePhoto, middleware))
	mux.Handle("PATCH /updateFirstName", AuthUser(userProfileHandler.UpdateFirstName, middleware))
	mux.Handle("PATCH /updateLastName", AuthUser(userProfileHandler.UpdateLastName, middleware))
	mux.Handle("PATCH /updateBio", AuthUser(userProfileHandler.UpdateBio, middleware))

	mux.HandleFunc("POST /ws/createRoom", wsHandler.CreateRoom)
	mux.HandleFunc("/ws/joinRoom/{roomId}", wsHandler.JoinRoom)
	mux.HandleFunc("GET /ws/getRooms", wsHandler.GetRooms)
	mux.HandleFunc("GET /ws/getClients", wsHandler.GetClients)

	// nantian
	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return mux
}

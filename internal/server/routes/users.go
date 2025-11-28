package routes

import (
	"encoding/json"
	"gohome/internal/core"
	"gohome/internal/repository"
	"gohome/internal/security"
	"gohome/shared/config"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type UsersRouter struct {
	store          *sessions.CookieStore
	userRepository *repository.UserRepository
	appEnv         config.AppEnv
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUsersRouter(mux *http.ServeMux, sessionSecret string, appEnv config.AppEnv, userRepository *repository.UserRepository) *UsersRouter {
	var store = sessions.NewCookieStore([]byte(sessionSecret))
	r := &UsersRouter{
		store:          store,
		userRepository: userRepository,
		appEnv:         appEnv,
	}

	mux.HandleFunc("POST /api/users/login", r.handleLogin)
	mux.HandleFunc("POST /api/users/logout", r.handleLogout)
	mux.HandleFunc("POST /api/users/register", r.handleRegister)
	mux.HandleFunc("GET /api/users/me", r.handleMe)
	mux.HandleFunc("GET /api/users/can_register", r.handleCanRegister)

	return r
}

func (s *UsersRouter) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest UserRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !mailIsValid(loginRequest.Email) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := s.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	passOk := security.CheckPasswordHash(loginRequest.Password, user.PasswordHash)
	if !passOk {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	session, _ := s.store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID

	session.Options.HttpOnly = true
	session.Options.Secure = s.appEnv == config.Production
	session.Options.SameSite = http.SameSiteStrictMode
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
}

func (s *UsersRouter) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session-name")
	session.Values = make(map[any]any)
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Error while logging out", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out"}`))
}

func (s *UsersRouter) handleCanRegister(w http.ResponseWriter, r *http.Request) {
	count, err := s.userRepository.Count()
	if err != nil || *count != 0 {
		json.NewEncoder(w).Encode(map[string]any{"can_register": false})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"can_register": true})
}

func (s *UsersRouter) handleMe(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (s *UsersRouter) handleRegister(w http.ResponseWriter, r *http.Request) {
	count, err := s.userRepository.Count()
	if err != nil || *count != 0 {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var signInRequest UserRequest

	if err := json.NewDecoder(r.Body).Decode(&signInRequest); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !mailIsValid(signInRequest.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := security.HashPassword(signInRequest.Password)
	if err != nil {
		http.Error(w, "Invalid password format", http.StatusBadRequest)
		return
	}

	u := &core.User{
		ID:           uuid.New().String(),
		Email:        signInRequest.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepository.Save(u); err != nil {
		http.Error(w, "error while saving user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (s *UsersRouter) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session-name")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func mailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

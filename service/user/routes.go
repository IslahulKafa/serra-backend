package user

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"serra/types"
	"serra/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/verify-otp", h.handleVerifyOTP).Methods("POST")
	router.HandleFunc("/refresh-token", h.handleRefreshToken).Methods("POST")
	router.Handle("/onboarding", utils.JWTAuth(http.HandlerFunc(h.handleOnboarding))).Methods("POST")
	router.Handle("/me", utils.JWTAuth(http.HandlerFunc(h.handleProfile))).Methods("GET")
	router.Handle("/keys/upload", utils.JWTAuth(http.HandlerFunc(h.handleUploadKeys))).Methods("POST")
	router.HandleFunc("/keys/{user_id}", h.handleGetPrekeyBundle).Methods("GET")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=20"`
	}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user := &types.User{
		Email:    payload.Email,
		Password: string(hashed),
	}

	if err := h.store.CreateUser(user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "Account succesfully created!",
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=20"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	code := fmt.Sprintf("%06d", rand.IntN(1000000))
	otpToken, err := utils.GenerateOTPToken(payload.Email, code, 5*time.Minute)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":   "Login successful!",
		"otp":       code,
		"otp_token": otpToken,
	})
}

func (h *Handler) handleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Code     string `json:"code" validate:"required,len=6"`
		OTPToken string `json:"otp_token" validate:"required"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	emailFromToken, otpFromToken, err := utils.VerifyOTPToken(payload.OTPToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if payload.Email != emailFromToken || payload.Code != otpFromToken {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid OTP"))
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("user not found"))
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	h.store.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(7*24*time.Hour))

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":       "OTP verified successfully",
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := h.store.GetRefreshToken(payload.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	token, err := utils.GenerateJWT(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
	})
}

func (h *Handler) handleOnboarding(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int64)

	var payload struct {
		Username   string `json:"username" validate:"required,min=3"`
		ProfilePic string `json:"profile_pic"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.ProfilePic == "" {
		payload.ProfilePic = "https://i.pinimg.com/736x/f9/24/12/f924127a8033eecbb67b0e1509097095.jpg"
	}

	err := h.store.SetUserProfile(userID, payload.Username, payload.ProfilePic)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Profile updated!",
	})
}

func (h *Handler) handleProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int64)

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func (h *Handler) handleUploadKeys(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int64)

	var payload struct {
		IdentityKey           string   `json:"identity_key" validate:"required"`
		SignedPrekey          string   `json:"signed_prekey" validate:"required"`
		SignedPrekeySignature string   `json:"signed_prekey_signature" validate:"required"`
		OneTimePrekeys        []string `json:"one_time_prekeys" validate:"required,dive,required"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.store.UpsertPrekeyBundle(userID, payload.IdentityKey, payload.SignedPrekey, payload.SignedPrekeySignature, payload.OneTimePrekeys)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Keys uploaded succesfully",
	})
}

func (h *Handler) handleGetPrekeyBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	bundle, err := h.store.GetPrekeyBundle(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, bundle)
}

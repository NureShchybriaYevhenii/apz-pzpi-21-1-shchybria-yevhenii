package handlers

import (
	"devquest-server/devquest/domain/models"
	"devquest-server/devquest/infrastructure"
	"devquest-server/devquest/usecases"
	"devquest-server/devquest/utils"
	"net/http"
)

type AuthHttpHandler struct {
	userUsecase usecases.UserUsecase
}

func NewAuthHttpHandler(userUsecase usecases.UserUsecase) *AuthHttpHandler {
	return &AuthHttpHandler{userUsecase: userUsecase}
}

func (a *AuthHttpHandler) Register(auth *infrastructure.Auth) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var userRegisterInfo models.RegisterUserDTO

		if err := utils.ReadJSON(w, r, &userRegisterInfo); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUserData, err := a.userUsecase.RegisterUser(userRegisterInfo)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUser := infrastructure.JWTUser{
			ID: jwtUserData.ID,
			Username: jwtUserData.Username,
			RoleTitle: jwtUserData.RoleTitle,
		}

		tokens, err := auth.GenerateTokenPairs(&jwtUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		refreshCookie := auth.GetRefreshCookie(tokens.RefreshToken)
		http.SetCookie(w, refreshCookie)

		utils.WriteJSON(w, http.StatusAccepted, tokens)
	}
}

func (a *AuthHttpHandler) Login(auth *infrastructure.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userLoginInfo models.LoginUserDTO

		if err := utils.ReadJSON(w, r, &userLoginInfo); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUserData, err := a.userUsecase.LoginUser(userLoginInfo)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUser := infrastructure.JWTUser{
			ID: jwtUserData.ID,
			Username: jwtUserData.Username,
			RoleTitle: jwtUserData.RoleTitle,
		}

		tokens, err := auth.GenerateTokenPairs(&jwtUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		refreshCookie := auth.GetRefreshCookie(tokens.RefreshToken)
		http.SetCookie(w, refreshCookie)

		utils.WriteJSON(w, http.StatusAccepted, tokens)
	}
}

func (a *AuthHttpHandler) Logout(auth *infrastructure.Auth) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, auth.GetExpiredRefreshCookie())
		w.WriteHeader(http.StatusAccepted)
	}
}
package handlers

import (
	"devquest-server/devquest/domain/models"
	"devquest-server/devquest/infrastructure"
	"devquest-server/devquest/usecases"
	"devquest-server/devquest/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHttpHandler struct {
	userUsecase usecases.UserUsecase
}

func NewUserHttpHandler(userUsecase usecases.UserUsecase) *UserHttpHandler {
	return &UserHttpHandler{userUsecase: userUsecase}
}

func (u *UserHttpHandler) Register(auth *infrastructure.Auth) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var userRegisterInfo models.RegisterUserDTO

		if err := utils.ReadJSON(w, r, &userRegisterInfo); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUserData, err := u.userUsecase.RegisterUser(userRegisterInfo)
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

		resData := struct {
			Tokens infrastructure.TokenPairs `json:"tokens"`
			UserID uuid.UUID `json:"user_id"`
			RoleTitle string `json:"role"`
		} {Tokens: tokens, UserID: jwtUser.ID, RoleTitle: jwtUser.RoleTitle}

		res := utils.JSONResponse{
			Error: false,
			Data: resData,
		}

		utils.WriteJSON(w, http.StatusAccepted, res)
	}
}

func (u *UserHttpHandler) Login(auth *infrastructure.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userLoginInfo models.LoginUserDTO

		if err := utils.ReadJSON(w, r, &userLoginInfo); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		jwtUserData, err := u.userUsecase.LoginUser(userLoginInfo)
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

		resData := struct {
			Tokens infrastructure.TokenPairs `json:"tokens"`
			UserID uuid.UUID `json:"user_id"`
			RoleTitle string `json:"role"`
		} {Tokens: tokens, UserID: jwtUser.ID, RoleTitle: jwtUser.RoleTitle}

		res := utils.JSONResponse{
			Error: false,
			Data: resData,
		}

		utils.WriteJSON(w, http.StatusAccepted, res)
	}
}

func (u *UserHttpHandler) Logout(auth *infrastructure.Auth) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, auth.GetExpiredRefreshCookie())
		w.WriteHeader(http.StatusAccepted)
	}
}

func (u *UserHttpHandler) GetDevelopersForManager(w http.ResponseWriter, r *http.Request) {
	managerID, err := uuid.Parse(chi.URLParam(r, "manager_id"))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	developers, err := u.userUsecase.GetDevelopersForManager(managerID)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, developers)
}

func (u *UserHttpHandler) GetRolesForRegistration(w http.ResponseWriter, r *http.Request) {
	roles, err := u.userUsecase.GetRolesForRegistration()
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, roles)
}
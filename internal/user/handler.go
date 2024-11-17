package user

import (
	"ChatApp/internal/image"
	"ChatApp/util"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"strings"
)

type UpdateDTO struct {
	ID             primitive.ObjectID `json:"id"`
	Username       string             `json:"username"`
	ProfilePicture *multipart.File    `json:"profile_picture"`
	Bio            *string            `json:"bio"`
	Email          *string            `json:"email"`
}

type UserHandler struct {
	usecase      Usecase
	imageUsecase image.Usecase
}

func NewUserHandler(usecase Usecase, imageUsecase image.Usecase) *UserHandler {
	return &UserHandler{usecase: usecase, imageUsecase: imageUsecase}
}

func (h *UserHandler) GetAllAvailableUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	tokenParts := strings.Split(token, " ")
	user, err := util.ParseToken(tokenParts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", user.ID)

	usersDTO, err := h.usecase.GetAllUsersDTO(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, usersDTO)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Fail to get user id", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.GetUserById(context.TODO(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Fail to get user id", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.GetUserById(context.TODO(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("profile_picture")
	if err != nil && err.Error() != "http: no such file" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file != nil {
		imageId, err := h.imageUsecase.CreateImage(context.TODO(), file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.ProfilePicture = &imageId
	}

	bio := r.PostFormValue("bio")
	user.Bio = &bio

	email := r.PostFormValue("email")
	user.Email = &email

	newToken, err := h.usecase.UpdateUser(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authCookie := http.Cookie{
		Name:  "auth_cookie",
		Value: newToken,
		Path:  "/",
	}

	http.SetCookie(w, &authCookie)

	util.JSONResponse(w, http.StatusOK, map[string]string{"access_token": newToken})
}

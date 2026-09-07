package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
	"match-me-api/internal/pkg/imgutil"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const maxProfilePictureSize = 1 << 20 // 1 MB

type UserManager interface {
	Profile(ctx context.Context, myID int64, targetUserID int64) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
	UpdateProfilePicture(ctx context.Context, userID int64, rawFile io.Reader) (string, error)
	DeleteProfilePicture(ctx context.Context, userID int64) (string, error)
}

type UserHandler struct {
	userService UserManager
	validator   *validator.Validate
	log         *slog.Logger
}

func NewUserHandler(um UserManager, v *validator.Validate, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: um, validator: v, log: logger}
}

// @UserSummary godoc
// @Summary      Get user summary
// @Description  Get user summary
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} dto.UserSummaryResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) UserSummary(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user id: NAN",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to get user",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.UserSummaryResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// @UserProfile godoc
// @Summary      Get user profile
// @Description  Get user profile with bio and age by user ID
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} dto.ProfileResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /users/{id}/profile [get]
func (h *UserHandler) UserProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user id: NAN",
		})
	}

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:  profile.UserID,
		Age: profile.Age,
		Bio: profile.Bio,
		// IsOnline: profile.IsOnline,
	})
}

// @UserBio godoc
// @Summary      Get user bio
// @Description  Get user biographical and activity data by user ID
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} dto.UserBioResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /users/{id}/bio [get]
func (h *UserHandler) UserBio(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user id: NAN",
		})
	}

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	activitiesResponse := make([]dto.UserActivityResponse, 0, len(profile.Activities))

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := dto.UserActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    int(a.Experience),
				InterestLevel: int(a.InterestLevel),
			}
			activitiesResponse = append(activitiesResponse, activity)
		}
	}

	return c.JSON(http.StatusOK, dto.UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: int(profile.InteractionMode),
		Activities:      activitiesResponse,
		Lat:             profile.Lat,
		Lon:             profile.Lon,
	})
}

// @MySummary godoc
// @Summary      Get current user summary
// @Description  Get authorized user summary with id, email, name, and picture URL
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.MySummaryResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me [get]
func (h *UserHandler) MySummary(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.MySummaryResponse{
		ID:         profile.UserID,
		Email:      profile.Email,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// @MyProfile godoc
// @Summary      Get current user profile
// @Description  Get authorized user profile with age and bio
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.ProfileResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/profile [get]
func (h *UserHandler) MyProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:  profile.UserID,
		Age: profile.Age,
		Bio: profile.Bio,
	})
}

// @MyBio godoc
// @Summary      Get current user bio
// @Description  Get authorized user bio and activities data
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.UserBioResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/bio [get]
func (h *UserHandler) MyBio(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	activitiesResponce := []dto.UserActivityResponse{}

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := dto.UserActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    int(a.Experience),
				InterestLevel: int(a.InterestLevel),
			}
			activitiesResponce = append(activitiesResponce, activity)
		}
	}

	return c.JSON(http.StatusOK, dto.UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: int(profile.InteractionMode),
		Activities:      activitiesResponce,
	})
}

// @UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update existing user profile fields
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ProfileUpdateRequest true "Profile update fields"
// @Success      200 {object} dto.OKResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/profile [patch]
func (h *UserHandler) UpdateProfile(c *echo.Context) error {
	const op = "httpserver.handlers.UpdateProfile"
	// log := h.log.With(slog.String("op", op))

	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	var req dto.ProfileUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	// 3. Map DTO into domain struct
	params := &domain.ProfileUpdateParams{
		Name:       req.Name,
		Age:        req.Age,
		PictureURL: req.PictureURL,
		Bio:        req.Bio,
		MaxRadius:  req.MaxRadius,
		Lat:        req.Lat,
		Lon:        req.Lon,
	}

	if req.InteractionMode != nil {
		params.InteractionMode = new(domain.InteractionMode(*req.InteractionMode))
	}

	if req.Activities != nil {
		activities := make([]domain.ActivityInput, len(*req.Activities))
		for i, act := range *req.Activities {
			activities[i] = domain.ActivityInput{
				ActivityID:    act.ID,
				Experience:    domain.ExperienceLevel(act.Experience),
				InterestLevel: domain.InterestLevel(act.InterestLevel),
			}
		}
		params.Activities = &activities
	}

	if err := h.userService.UpdateProfile(ctx, userID, params); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, dto.OKResponse{
		Message: "Profile updated successfully",
	})
}

// @UploadProfilePicture godoc
// @Summary      Upload profile picture
// @Description  Upload, process, and update the profile picture for the authorized user
// @Tags         user
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        picture formData file true "Profile picture file (max 1MB)"
// @Success      200 {object} dto.ProfilePictureUpdatedResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      413 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/picture [post]
func (h *UserHandler) UploadProfilePicture(c *echo.Context) error {
	const op = "httpserver.handlers.UploadProfilePicture"
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	fileHeader, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Form field 'picture' is required",
		})
	}

	if fileHeader.Size > maxProfilePictureSize {
		return c.JSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{
			Message: "File size exceeds 1MB limit",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Failed to read uploaded file",
		})
	}
	defer file.Close()

	pictureURL, err := h.userService.UpdateProfilePicture(ctx, userID, file)
	if err != nil {
		if errors.Is(err, imgutil.ErrUnsupportedFormat) || errors.Is(err, imgutil.ErrInvalidDimensions) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Invalid file format or dimensions",
			})
		}
		h.log.Error("failed to update profile picture", slog.String("op", op), "user_id", userID, logger.Err(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, dto.ProfilePictureUpdatedResponse{
		Message:    "Profile picture updated successfully",
		PictureURL: pictureURL,
	})
}

// @DeleteProfilePicture godoc
// @Summary      Delete profile picture
// @Description  Reset the profile picture for the authorized user to the default placeholder
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.ProfilePictureUpdatedResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/picture [delete]
func (h *UserHandler) DeleteProfilePicture(c *echo.Context) error {
	const op = "httpserver.handlers.DeleteProfilePicture"
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	defaultURL, err := h.userService.DeleteProfilePicture(ctx, userID)
	if err != nil {
		h.log.Error("failed to delete profile picture", slog.String("op", op), "user_id", userID, logger.Err(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, dto.ProfilePictureUpdatedResponse{
		Message:    "Profile picture removed successfully",
		PictureURL: defaultURL,
	})
}

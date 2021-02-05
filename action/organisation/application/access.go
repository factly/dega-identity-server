package application

import (
	"net/http"
	"strconv"

	"github.com/factly/kavach-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// access - Get access of application based on slug
// @Summary Get access of application based on slug
// @Description Get access of application based on slug
// @Tags OrganisationApplications
// @ID get-access-organisation-application
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param organisation_id path string true "Organisation ID"
// @Param application_slug path string true "Application Slug"
// @Success 200
// @Failure 401
// @Router /organisations/{organisation_id}/applications/{application_slug}/access [get]
func access(w http.ResponseWriter, r *http.Request) {
	appSlug := chi.URLParam(r, "application_slug")
	if appSlug == "" {
		errorx.Render(w, errorx.Parser(errorx.GetMessage("invalid slug", http.StatusBadRequest)))
		return
	}

	orgnaisationID := chi.URLParam(r, "organisation_id")
	oID, err := strconv.Atoi(orgnaisationID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	uID, err := strconv.Atoi(r.Header.Get("X-User"))
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	application := model.Application{}
	err = model.DB.Model(&model.Application{}).Where(&model.Application{
		OrganisationID: uint(oID),
		Slug:           appSlug,
	}).Preload("Users").First(&application).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	appUsers := application.Users

	for _, usr := range appUsers {
		if usr.ID == uint(uID) {
			loggerx.Error(err)
			renderx.JSON(w, http.StatusOK, nil)
			return
		}
	}

	renderx.JSON(w, http.StatusUnauthorized, nil)
}

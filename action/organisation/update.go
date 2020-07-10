package organisation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/kavach-server/model"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

func update(w http.ResponseWriter, r *http.Request) {

	req := &model.Organisation{}
	json.NewDecoder(r.Body).Decode(&req)

	organisationID := chi.URLParam(r, "organisation_id")
	orgID, err := strconv.Atoi(organisationID)

	organisation := &model.Organisation{}
	organisation.ID = uint(orgID)

	// check record exists or not
	err = model.DB.First(&organisation).Error
	if err != nil {
		return
	}

	// check the permission of host
	hostID, _ := strconv.Atoi(r.Header.Get("X-User"))
	permission := &model.OrganisationUser{}

	err = model.DB.Model(&model.OrganisationUser{}).Where(&model.OrganisationUser{
		OrganisationID: uint(orgID),
		UserID:         uint(hostID),
		Role:           "owner",
	}).First(permission).Error

	if err != nil {
		return
	}

	// delete
	model.DB.Model(&organisation).Updates(model.Organisation{
		Title: req.Title,
	}).First(&organisation)

	result := &orgWithRole{}
	result.Organisation = *organisation
	result.Permission = *permission

	renderx.JSON(w, http.StatusOK, result)
}
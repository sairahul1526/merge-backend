package admin

import (
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	"net/http"
	"strconv"

	UTIL "merge-backend/util"
)

// ItemAdd - add new item to stock
func ItemAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ItemAddAdminRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, "", response)
		return
	}

	// check if stock is positive
	stock, err := strconv.Atoi(body["stock"])
	if stock < 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.StockNegativeUnacceptedMessage, CONSTANT.ShowDialog, "", response)
		return
	}
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// create item
	itemID, _, err := DB.MainDB.InsertWithUniqueID(CONSTANT.ItemsTable, map[string]string{
		"title":       body["title"],
		"description": body["description"],
		"stock":       body["stock"],
		"created_by":  r.Header.Get("user_id"),
		"status":      CONSTANT.ItemActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	response["item_id"] = itemID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

// ItemUpdate - update item, maintain stock
func ItemUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ItemUpdateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, "", response)
		return
	}

	// check if stock is positive
	stock, err := strconv.Atoi(body["stock"])
	if stock < 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.StockNegativeUnacceptedMessage, CONSTANT.ShowDialog, "", response)
		return
	}
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// update item
	_, err = DB.MainDB.UpdateSQL(CONSTANT.ItemsTable, map[string]string{
		"id": r.FormValue("item_id"),
	}, map[string]string{
		"title":       body["title"],
		"description": body["description"],
		"stock":       body["stock"],
		"updated_by":  r.Header.Get("user_id"),
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

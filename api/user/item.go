package user

import (
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	"net/http"
	"strconv"

	UTIL "merge-backend/util"
)

// ItemsGet - get all available items
func ItemsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	where := " where stock > 0 and status = " + CONSTANT.ItemActive + " "
	// get available items
	items, err := DB.MainDB.SelectProcess("select id, title, description, stock, status from " + CONSTANT.ItemsTable + where + " order by created_at desc limit " + strconv.Itoa(CONSTANT.ResultsPerPageUser) + " offset " + strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageUser))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// get total number of items
	itemsCount, err := DB.MainDB.SelectProcess("select count(*) as ctn from " + CONSTANT.ItemsTable + where)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	response["items"] = items
	response["items_count"] = itemsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(itemsCount[0]["ctn"], CONSTANT.ResultsPerPageUser))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

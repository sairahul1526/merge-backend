package user

import (
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	"net/http"
	"strings"

	UTIL "merge-backend/util"
)

// CartUpdate - add/remove in cart
func CartUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.CartUpdateUserRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, "", response)
		return
	}

	if strings.EqualFold(body["type"], "add") {

		// get item count in cart
		count, _ := DB.MainDB.QueryRowSQL("select count from "+CONSTANT.CartsTable+" where user_id = $1 and item_id = $2", r.Header.Get("user_id"), body["item_id"])
		if len(count) == 0 {
			count = "0"
		}

		// check if stock available
		items, err := DB.MainDB.SelectProcess("select 1 from "+CONSTANT.ItemsTable+" where id = $1 and stock > $2", body["item_id"], count)
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
			return
		}
		if len(items) == 0 {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.StockUnavailableMessage, CONSTANT.ShowDialog, "", response)
			return
		}

		// add to cart if not available
		// if available, increase count
		_, err = DB.MainDB.ExecuteSQL("insert into "+CONSTANT.CartsTable+" (user_id, item_id, count) values ($1, $2, 1) on conflict (user_id, item_id) do update set count = "+CONSTANT.CartsTable+".count + 1", r.Header.Get("user_id"), body["item_id"])
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
			return
		}
	} else {
		// decrease count
		_, err = DB.MainDB.ExecuteSQL("update "+CONSTANT.CartsTable+" set count = count - 1 where user_id = $1 and item_id = $2", r.Header.Get("user_id"), body["item_id"])
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
			return
		}

		// delete if item count is zero
		DB.MainDB.DeleteSQL(CONSTANT.CartsTable, map[string]string{
			"user_id": r.Header.Get("user_id"),
			"item_id": body["item_id"],
			"count":   "0",
		})
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

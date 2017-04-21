package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/src/models"
	"github.com/mamoroom/echo-mvc/src/router/handler"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	"errors"
	"fmt"
	_ "reflect"
)

type ReqOpening struct {
	Name string `json:"name" validate:"required"`
}

func OpeningHandler(c echo.Context) error {
	//必須
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "OpeningHandler"

	var req_open = new(ReqOpening)
	if err := c.Bind(req_open); err != nil {
		return res_json.ErrorBadRequest(c, "InvalidRequest", err, "Could not bind request body")
	}
	if err := c.Validate(req_open); err != nil {
		return res_json.Failed(c, "ValidationFailure")
	}

	////// UserW //////
	user_w := models.NewUserW()
	user_w.Dbh.SetNewSession()
	rollback_func := func() error { return user_w.Dbh.Rollback() }
	commit_func := func() error { return user_w.Dbh.Commit() }
	defer user_w.Dbh.Close()
	defer fmt.Println("End Transaction.")

	// Transaction //
	if err := user_w.Dbh.BeginTx(); err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not begin transaction")
	}

	user_w.Dbh.ForUpdate()
	if err := user_w.FindUserById(res_jwt.Data.SessionUserModel.GetUserEntity().Id); err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not get user data from master DB")
	}
	rows_affected, err := user_w.UpdateTutorialDone(req_open.Name)
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not insert user data")
	}
	if rows_affected == 0 {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on insert auth data"), "Could not insert data")
	}
	handle_rollback_or_commit(commit_func)
	///////////////////

	res_jwt.SetSessionUser(user_w)
	return NazoHandler(c)
}

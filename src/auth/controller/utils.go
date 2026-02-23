package auth

//func SendJSONErrorResponse(ctx echo.Context, err error) error {
//	var badRequestError *common.BadRequestError
//	var unauthorizedError *common.UnauthorizedError
//	var notFoundError *common.NotFoundError
//	var conflictError *common.ConflictError
//	var unprocessableError *common.UnprocessableError
//
//	apiCode := http.StatusInternalServerError
//
//	switch {
//	case errors.As(err, &badRequestError):
//		apiCode = http.StatusBadRequest
//	case errors.As(err, &unauthorizedError):
//		apiCode = http.StatusUnauthorized
//	case errors.As(err, &notFoundError):
//		apiCode = http.StatusNotFound
//	case errors.As(err, &conflictError):
//		apiCode = http.StatusConflict
//	case errors.As(err, &unprocessableError):
//		apiCode = http.StatusUnprocessableEntity
//	}
//
//	apiError := convertToError(err)
//
//	return ctx.JSON(apiCode, apiError)
//}
//
//func convertToError(err error) api.Error {
//	var appError common.AppError
//	if errors.As(err, &appError) {
//		return api.Error{
//			Code:    appError.GetCode(),
//			Message: appError.GetMessage(),
//		}
//	}
//
//	return api.Error{
//		Code:    "UNEXPECTED",
//		Message: err.Error(),
//	}
//}
//
//func GetUserId(ctx echo.Context) (userId int, err error) {
//	cookie, err := ctx.Cookie(SessionCookieName)
//
//	if err != nil {
//		err = common.NewCookieNotFoundUnauthorizedError()
//		return
//	}
//
//	sessionId := cookie.Value
//
//	userSession, err := auth_service.GetUserSessionById(sessionId)
//
//	if errors.Is(err, sql.ErrNoRows) {
//		err = common.NewActiveSessionNotFoundUnauthorizedError()
//		return
//	}
//
//	if err != nil {
//		return
//	}
//
//	userId = userSession.UserId
//
//	return
//}
//
//func doesUserSessionExist(ctx context.Context) (doesExist bool, err error) {
//	cookie, err := getSessionCookie(ctx)
//
//	if err != nil {
//		err = common.NewCookieNotFoundUnauthorizedError()
//		return
//	}
//
//	sessionId := cookie.Value
//	_, err = auth_service.GetUserSessionById(sessionId)
//
//	if errors.Is(err, sql.ErrNoRows) {
//		err = nil
//		return
//	}
//
//	if err != nil {
//		return
//	}
//
//	doesExist = true
//	return
//}
//
//func getSessionCookie(ctx context.Context) (*http.Cookie, error) {
//	cookie, err := ctx.Cookie(SessionCookieName)
//
//	if err != nil {
//		return nil, common.NewCookieNotFoundUnauthorizedError()
//	}
//
//	return cookie, nil
//}

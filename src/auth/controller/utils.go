package auth

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

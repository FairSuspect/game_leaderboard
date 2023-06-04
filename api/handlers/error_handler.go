package handlers

import (
	"game_leaderboards/m/v2/api/errors"
	"net/http"
)

// Обработчик ошибок
func HandleError(err error, w http.ResponseWriter) {
	switch e := err.(type) {
	case errors.NotFoundError:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(e.Error()))
		return
	case errors.DecodeError:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
		return

	default:

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		// Обработка других ошибок
	}
}

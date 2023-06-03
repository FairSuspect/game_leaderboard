package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ParseID(idRaw string) (int, error) {
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ParseIdFromRequest(r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	idRaw := vars["id"]

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		fmt.Println("Ошибка преобразования строки в int:", err)
		return 0, true
	}
	return id, false
}

package handlers

import (
	"FetchChallenge/rules"
	"FetchChallenge/store"
	"FetchChallenge/types"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func writeReceiptNotFoundError(w http.ResponseWriter) {
	http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
}

func writePostReceiptError(w http.ResponseWriter) {
	http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This route is invalid.", http.StatusNotFound)
}

func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		writeReceiptNotFoundError(w)
		return
	}
	if !uuidRegex.MatchString(id) {
		writeReceiptNotFoundError(w)
		return
	}

	points, ok := store.ReceiptStore.Get(id)
	if !ok {
		writeReceiptNotFoundError(w)
		return
	}

	json, _ := json.Marshal(types.Score{Points: strconv.FormatInt(points, 10)})
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func PostReceiptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writePostReceiptError(w)
		return
	}

	var receipt types.Receipt
	if err := json.Unmarshal(body, &receipt); err != nil {
		writePostReceiptError(w)
		return
	}

	newId, err := uuid.NewRandom()
	newIdString := newId.String()
	if err != nil {
		writePostReceiptError(w)
		return
	}

	var total int64 = 0
	for _, rule := range rules.ReceiptRules {
		total += rule(&receipt)
	}

	store.ReceiptStore.Set(newIdString, total)

	json, _ := json.Marshal(types.ReceiptCreated{Id: newIdString})
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

package queue

import (
	"encoding/json"
	"net/http"

	"github.com/bigspaceships/circlejerk/auth"
)

type QueueEntry struct {
	Name     string `json: "name"`
	Username string `json: "username"`
	Type     string `json: "type"`
}

var queue []QueueEntry

func SetupQueue() {
	queue = make([]QueueEntry, 0)
}

func JoinQueue(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     "Clarifier",
	}

	queue = append(queue, newEntry)

	w.WriteHeader(http.StatusOK)
}

func GetQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(queue)
}

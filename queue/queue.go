package queue

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"github.com/bigspaceships/circlejerk/auth"
	dq_websocket "github.com/bigspaceships/circlejerk/websocket"
)

type QueueEntry struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Id       int    `json:"id"`
}

type QueueRequestData struct {
	Type string `json:"type"`
}

type Queue struct {
	Points     []QueueEntry `json:"points"`
	Clarifiers []QueueEntry `json:"clarifiers"`
	Topic      string       `json:"topic"`
	wsServer   *dq_websocket.WsServer
	pointCount int
}

func SetupQueue(wsServer *dq_websocket.WsServer) *Queue {
	return &Queue{
		Points:     make([]QueueEntry, 0),
		Clarifiers: make([]QueueEntry, 0),
		Topic:      "Big long discussion",
		wsServer:   wsServer,
		pointCount: 0,
	}
}

func (queue *Queue) DeletePoint(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "error parsing id to int"+err.Error(), http.StatusBadRequest)
		return
	}

	pointIndex := slices.IndexFunc(queue.Points, func(entry QueueEntry) bool {
		return entry.Id == id
	})

	point := queue.Points[pointIndex]

	if !(userInfo.IsEboard || point.Username == userInfo.Username) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	queue.Points = slices.Concat(queue.Points[:pointIndex], queue.Points[(pointIndex+1):])
	if queue.Points == nil {
		queue.Points = make([]QueueEntry, 0)
	}

	queue.wsServer.SendWSMessage(struct {
		Type      string `json:"type"`
		Id        int    `json:"id"`
		Dismisser string `json:"dismisser"`
	}{
		Type:      "delete",
		Id:        id,
		Dismisser: userInfo.Name,
	})
}

func (queue *Queue) DeleteClarifier(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "error parsing id to int"+err.Error(), http.StatusBadRequest)
		return
	}

	pointIndex := slices.IndexFunc(queue.Clarifiers, func(entry QueueEntry) bool {
		return entry.Id == id
	})

	point := queue.Clarifiers[pointIndex]

	if !(userInfo.IsEboard || point.Username == userInfo.Username) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	queue.Clarifiers = slices.Concat(queue.Clarifiers[:pointIndex], queue.Clarifiers[(pointIndex+1):])
	if queue.Clarifiers == nil {
		queue.Clarifiers = make([]QueueEntry, 0)
	}

	queue.wsServer.SendWSMessage(struct {
		Type      string `json:"type"`
		Id        int    `json:"id"`
		Dismisser string `json:"dismisser"`
	}{
		Type:      "delete",
		Id:        id,
		Dismisser: userInfo.Name,
	})
}

func (queue *Queue) NewClarifier(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	queue.pointCount++

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     "clarifier",
		Id:       queue.pointCount,
	}

	queue.Clarifiers = append(queue.Clarifiers, newEntry)

	w.WriteHeader(http.StatusOK)

	queue.wsServer.SendWSMessage(struct {
		Type string     `json:"type"`
		Data QueueEntry `json:"data"`
	}{
		Type: "clarifier",
		Data: newEntry,
	})
}

func (queue *Queue) NewPoint(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	queue.pointCount++

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     "point",
		Id:       queue.pointCount,
	}

	queue.Points = append(queue.Points, newEntry)

	w.WriteHeader(http.StatusOK)

	queue.wsServer.SendWSMessage(struct {
		Type string     `json:"type"`
		Data QueueEntry `json:"data"`
	}{
		Type: "point",
		Data: newEntry,
	})
}

func (queue *Queue) ChangeTopic(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	if !userInfo.IsEboard {
		http.Error(w, "user is not on eboard", http.StatusForbidden)
		return
	}

	requestData := struct {
		NewTopic string `json:"new-topic"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		http.Error(w, "Error decoding body:"+err.Error(), http.StatusBadRequest)
		return
	}

	queue.Topic = requestData.NewTopic

	queue.wsServer.SendWSMessage(struct {
		Type  string `json:"type"`
		Topic string `json:"topic"`
	}{
		Type:  "topic",
		Topic: requestData.NewTopic,
	})

	w.WriteHeader(http.StatusOK)
}

func (queue *Queue) GetQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(queue)
}

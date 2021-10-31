package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"deduplicator/internal/models"
	"deduplicator/internal/services"
	"deduplicator/internal/services/mappers"

	"cloud.google.com/go/pubsub"
)

const (
	nextHopMethodKey = "next_hop_method"
	nextHopKey       = "next_hop"
	pubsubMethod     = "pubsub"
	httpMethod       = "http"
	hashAttributeKey = "hash_"
)

func NewDeduplicatorController(deduplicatorService services.IDeduplicatorService,
	pubSubService *services.PubSubService,
	httpService services.IHTTPService) *DeduplicatorController {
	return &DeduplicatorController{deduplicatorService: deduplicatorService, pubSubService: pubSubService, httpService: httpService}
}

type DeduplicatorController struct {
	pubSubService       *services.PubSubService
	deduplicatorService services.IDeduplicatorService
	httpService         services.IHTTPService
}

func (controller *DeduplicatorController) HandlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")

	request, err := mappers.ReaderToRequest(r.Body)
	if err != nil {
		controller.writeHTTPErrorAndLog(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mustForward, err := controller.deduplicatorService.CreateOrUpdate(controller.getDeduplicatorAttributes(request.Message.Attributes))
	if err != nil {
		controller.writeHTTPErrorAndLog(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !mustForward {
		return
	}

	code, err := controller.forwardMessage(request)
	if err != nil {
		controller.writeHTTPErrorAndLog(w, err.Error(), code)
		return
	}
}

func (controller *DeduplicatorController) writeHTTPErrorAndLog(w http.ResponseWriter, message string, code int) {
	log.Println(message)
	http.Error(w, message, code)
}

func (controller *DeduplicatorController) getDeduplicatorAttributes(attributes map[string]string) []string {
	var deduplicatorAttributes []string
	mappedHash := make(map[string]string)
	for key := range attributes {
		if strings.Contains(key, hashAttributeKey) {
			splittedHash := strings.Split(key, "_")
			mappedHash[splittedHash[1]] = key
		}
	}
	// this loop begins at 1 to be able to find every attributes to use to deduplicate with key from hash_1 to hash_n
	for i := 1; i < len(mappedHash)+1; i++ {
		deduplicatorAttributes = append(deduplicatorAttributes, attributes[mappedHash[strconv.Itoa(i)]])
	}
	return deduplicatorAttributes
}

func (controller *DeduplicatorController) deleteDeduplicatorAttributes(attributes map[string]string) {
	delete(attributes, nextHopKey)
	delete(attributes, nextHopMethodKey)
	for key := range attributes {
		if strings.Contains(key, hashAttributeKey) {
			delete(attributes, key)
		}
	}
}

func (controller *DeduplicatorController) forwardMessage(request *models.Request) (int, error) {
	destination := request.Message.Attributes[nextHopKey]
	log.Printf("destination: %s", destination)
	if destination == "" {
		return http.StatusBadRequest, errors.New("no destination given")
	}
	method := request.Message.Attributes[nextHopMethodKey]
	if method != pubsubMethod && method != httpMethod {
		return http.StatusBadRequest, errors.New("method must be http or pubsub")
	}
	controller.deleteDeduplicatorAttributes(request.Message.Attributes)
	if method == pubsubMethod {
		pubSubMessage := pubsub.Message{Data: request.Message.Data, Attributes: request.Message.Attributes}
		err := controller.pubSubService.Publish(destination, &pubSubMessage)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		log.Printf("message published to PubSub topic. topic path is: %s ", destination)
		return http.StatusOK, nil
	}
	err := controller.httpService.Post(destination, request.Message.Data, request.Message.Attributes)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	log.Printf("message forwarded by http request, destination url is: %s", destination)
	return http.StatusOK, nil
}

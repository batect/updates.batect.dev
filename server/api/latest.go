// Copyright 2019-2021 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

package api

import (
	"net/http"

	"github.com/batect/service-observability/middleware"
	"github.com/batect/updates.batect.dev/server/events"
	"github.com/batect/updates.batect.dev/server/storage"
)

type latestHandler struct {
	store     storage.LatestVersionStore
	eventSink events.EventSink
}

func NewLatestHandler(store storage.LatestVersionStore, eventSink events.EventSink) http.Handler {
	return &latestHandler{
		store:     store,
		eventSink: eventSink,
	}
}

func (h *latestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !requireMethod(w, req, http.MethodGet) {
		return
	}

	log := middleware.LoggerFromContext(req.Context())

	descriptor, err := h.store.GetLatestVersionDescriptor(req.Context())

	if err != nil {
		log.WithError(err).Error("Getting latest version descriptor failed.")
		serviceUnavailable(req.Context(), w)

		return
	}

	h.eventSink.PostLatestVersionCheck(req.Context(), req.UserAgent())

	w.Header().Set(contentTypeHeader, descriptor.ContentType)

	if _, err := w.Write(descriptor.Content); err != nil {
		log.WithError(err).Error("Writing response failed.")
		return
	}
}

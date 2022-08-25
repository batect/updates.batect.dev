// Copyright 2019-2022 Charles Korn.
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
	"regexp"

	"github.com/batect/updates.batect.dev/server/events"
)

type filesHandler struct {
	urlPattern *regexp.Regexp
	eventSink  events.EventSink
}

func NewFilesHandler(eventSink events.EventSink) http.Handler {
	return &filesHandler{
		urlPattern: regexp.MustCompile(`^/v1/files/(?P<versionInPath>\d+\.\d+\.\d+)/batect-(?P<versionInFileName>\d+\.\d+\.\d+).jar$`),
		eventSink:  eventSink,
	}
}

func (h *filesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//nolint:contextcheck
	if !requireMethod(w, req, http.MethodGet) {
		return
	}

	match := h.urlPattern.FindStringSubmatch(req.URL.Path)

	if match == nil {
		http.NotFound(w, req)
		return
	}

	versionInPath := match[1]

	if versionInFileName := match[2]; versionInPath != versionInFileName {
		http.NotFound(w, req)
		return
	}

	version := versionInPath
	fileName := "batect-" + version + ".jar"

	h.eventSink.PostFileDownload(req.Context(), req.UserAgent(), version, fileName)

	w.Header().Set("Location", "https://github.com/batect/batect/releases/download/"+versionInPath+"/"+fileName)
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.WriteHeader(http.StatusFound)
}

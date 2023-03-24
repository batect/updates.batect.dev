// Copyright 2019-2023 Charles Korn.
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

package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/batect/services-common/middleware/testutils"
	"github.com/batect/updates.batect.dev/server/api"
	"github.com/batect/updates.batect.dev/server/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Latest version endpoint", func() {
	var eventSink *mockEventSink
	var handler http.Handler
	var latestVersionStoreMock *mockLatestVersionStore
	var resp *httptest.ResponseRecorder

	BeforeEach(func() {
		eventSink = newMockEventSink()
		latestVersionStoreMock = &mockLatestVersionStore{}
		handler = api.NewLatestHandler(latestVersionStoreMock, eventSink)
		resp = httptest.NewRecorder()
	})

	Context("when invoked with a HTTP method other than GET", func() {
		BeforeEach(func() {
			req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("POST", "/v1/latest", nil))
			handler.ServeHTTP(resp, req)
		})

		It("returns a HTTP 405 response", func() {
			Expect(resp.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("returns a JSON error payload", func() {
			Expect(resp.Body).To(MatchJSON(`{"message":"This endpoint only supports GET requests"}`))
		})

		It("sets the response Content-Type header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("sets the response Allow header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Allow", []string{"GET"}))
		})

		It("does not post any events", func() {
			Expect(eventSink.LatestVersionCheckEventsPosted).To(BeEmpty())
		})
	})

	Context("when invoked with a HTTP GET", func() {
		var req *http.Request

		BeforeEach(func() {
			req, _ = testutils.RequestWithTestLogger(httptest.NewRequest("GET", "/v1/latest", nil))
			req.Header.Set("User-Agent", "MyApp/1.2.3")
		})

		Context("given retrieving the latest version information succeeds", func() {
			BeforeEach(func() {
				latestVersionStoreMock.errorToReturn = nil
				latestVersionStoreMock.descriptorToReturn = storage.VersionDescriptor{
					Content:     []byte(`{"some":"descriptor"}`),
					ContentType: "application/json+descriptor",
				}

				handler.ServeHTTP(resp, req)
			})

			It("returns a HTTP 200 response", func() {
				Expect(resp.Code).To(Equal(http.StatusOK))
			})

			It("returns the version descriptor in the response body", func() {
				Expect(resp.Body.String()).To(Equal(`{"some":"descriptor"}`))
			})

			It("returns the content type provided by the version information source", func() {
				Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json+descriptor"}))
			})

			It("posts a 'latest version check' event", func() {
				Expect(eventSink.LatestVersionCheckEventsPosted).To(ConsistOf(latestVersionCheckEvent{
					userAgent: "MyApp/1.2.3",
				}))
			})
		})

		Context("given retrieving the latest version information fails", func() {
			BeforeEach(func() {
				latestVersionStoreMock.errorToReturn = errors.New("something went wrong")
				latestVersionStoreMock.descriptorToReturn = storage.VersionDescriptor{}

				handler.ServeHTTP(resp, req)
			})

			It("returns a HTTP 503 response", func() {
				Expect(resp.Code).To(Equal(http.StatusServiceUnavailable))
			})

			It("returns a JSON error payload", func() {
				Expect(resp.Body).To(MatchJSON(`{"message":"Service unavailable"}`))
			})

			It("sets the response Content-Type header", func() {
				Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			})

			It("does not post any events", func() {
				Expect(eventSink.LatestVersionCheckEventsPosted).To(BeEmpty())
			})
		})
	})
})

type mockLatestVersionStore struct {
	descriptorToReturn storage.VersionDescriptor
	errorToReturn      error
}

func (m *mockLatestVersionStore) GetLatestVersionDescriptor(_ context.Context) (storage.VersionDescriptor, error) {
	return m.descriptorToReturn, m.errorToReturn
}

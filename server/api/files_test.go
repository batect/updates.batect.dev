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
	"net/http"
	"net/http/httptest"

	"github.com/batect/services-common/middleware/testutils"
	"github.com/batect/updates.batect.dev/server/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Files endpoint", func() {
	var eventSink *mockEventSink
	var handler http.Handler
	var resp *httptest.ResponseRecorder

	BeforeEach(func() {
		eventSink = newMockEventSink()
		handler = api.NewFilesHandler(eventSink)
		resp = httptest.NewRecorder()
	})

	Context("when invoked with a HTTP method other than GET", func() {
		BeforeEach(func() {
			req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("POST", "/v1/files/0.1.2/batect-0.1.2.jar", nil))
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
			Expect(eventSink.FileDownloadEventsPosted).To(BeEmpty())
		})
	})

	Context("when invoked with a HTTP GET", func() {
		Context("when invoked with a valid path", func() {
			BeforeEach(func() {
				req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("GET", "/v1/files/0.1.2/batect-0.1.2.jar", nil))
				req.Header.Set("User-Agent", "MyApp/1.2.3")

				handler.ServeHTTP(resp, req)
			})

			It("returns a HTTP 302 response", func() {
				Expect(resp.Code).To(Equal(http.StatusFound))
			})

			It("returns the GitHub download URL in the Location header", func() {
				Expect(resp.Header()).To(HaveKeyWithValue("Location", []string{"https://github.com/batect/batect/releases/download/0.1.2/batect-0.1.2.jar"}))
			})

			It("does not return a response body", func() {
				Expect(resp.Body.String()).To(BeEmpty())
			})

			It("does not set the response Content-Type header", func() {
				Expect(resp.Result().Header).ToNot(HaveKey("Content-Type"))
			})

			It("prevents caching of the response", func() {
				Expect(resp.Result().Header).To(HaveKeyWithValue("Cache-Control", []string{"no-store, max-age=0"}))
			})

			It("posts a 'file download' event", func() {
				Expect(eventSink.FileDownloadEventsPosted).To(ConsistOf(fileDownloadEvent{
					userAgent: "MyApp/1.2.3",
					version:   "0.1.2",
					fileName:  "batect-0.1.2.jar",
				}))
			})
		})

		Context("when invoked with an invalid path", func() {
			examples := []string{
				"/",
				"/v1",
				"/v1/files",
				"/v1/files/",
				"/v1/files/0.1.2",
				"/v1/files/0.1.2/",
				"/v1/files/0.1.2/batect.jar",
				"/v1/files/0.1.2/batect-0.1.2",
				"/v1/files/0.1.2/batect-0.1.2.blah",
				"/v1/files/0.1/batect-0.1.2.jar",
				"/v1/files/0/batect-0.1.2.jar",
				"/v1/files/blah/batect-0.1.2.jar",
				"/v1/files/0.1.2/batect-0.1.jar",
				"/v1/files/0.1.2/batect-0.jar",
				"/v1/files/0.1.2/batect-blah.jar",
				"/v1/files/0.1.2/batect-3.4.5.jar",
				"/v1/files/0.1.2/batect-0.1.2.jar/thing",
				"/v1/files/0.1.2/somethingelse-0.1.2.jar",
			}

			for _, e := range examples {
				path := e

				Context("given the invalid path '"+path+"'", func() {
					BeforeEach(func() {
						req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("GET", path, nil))
						handler.ServeHTTP(resp, req)
					})

					It("returns a HTTP 404 response", func() {
						Expect(resp.Code).To(Equal(http.StatusNotFound))
					})

					It("includes the default Golang 404 error message in the body", func() {
						Expect(resp.Body.String()).To(Equal("404 page not found\n"))
					})

					It("does not post any events", func() {
						Expect(eventSink.FileDownloadEventsPosted).To(BeEmpty())
					})
				})
			}
		})

	})
})

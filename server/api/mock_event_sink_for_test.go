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

import "context"

type mockEventSink struct {
	LatestVersionCheckEventsPosted []latestVersionCheckEvent
	FileDownloadEventsPosted       []fileDownloadEvent
}

type latestVersionCheckEvent struct {
	userAgent string
}

type fileDownloadEvent struct {
	userAgent string
	version   string
	fileName  string
}

func newMockEventSink() *mockEventSink {
	return &mockEventSink{
		LatestVersionCheckEventsPosted: []latestVersionCheckEvent{},
		FileDownloadEventsPosted:       []fileDownloadEvent{},
	}
}

func (m *mockEventSink) PostLatestVersionCheck(_ context.Context, userAgent string) {
	m.LatestVersionCheckEventsPosted = append(
		m.LatestVersionCheckEventsPosted,
		latestVersionCheckEvent{userAgent: userAgent},
	)
}

func (m *mockEventSink) PostFileDownload(_ context.Context, userAgent string, version string, fileName string) {
	m.FileDownloadEventsPosted = append(
		m.FileDownloadEventsPosted,
		fileDownloadEvent{
			userAgent: userAgent,
			version:   version,
			fileName:  fileName,
		},
	)
}

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

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type test interface {
	Description() string
	Run(baseUrl string) error
}

func main() {
	baseUrl, err := getBaseUrl()

	if err != nil {
		fmt.Printf("Getting base URL failed: %s\n", err)
		os.Exit(1)
	}

	tests := []test{
		&pingTest{},
		&latestTest{},
		&downloadTest{},
	}

	for _, t := range tests {
		fmt.Printf("Running: %s\n", t.Description())

		if err := t.Run(baseUrl); err != nil {
			fmt.Printf("> Test failed!\n")
			fmt.Printf("> %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("> Test passed.\n")
		fmt.Println()
	}

	fmt.Println("All tests passed.")
}

func getBaseUrl() (string, error) {
	domain := os.Getenv("DOMAIN")

	if domain == "" {
		return "", errors.New("environment variable 'DOMAIN' is not set")
	}

	return fmt.Sprintf("https://%s", domain), nil
}

func makeRequest(baseUrl string, path string) (*http.Response, error) {
	clientWithNoRedirectFollowing := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("User-Agent", "UpdatesServiceSmokeTest/2.0.0")

	return clientWithNoRedirectFollowing.Do(req)
}

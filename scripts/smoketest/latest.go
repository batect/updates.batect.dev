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
	"encoding/json"
	"fmt"
	"strings"
)

type latestTest struct{}

func (t *latestTest) Description() string {
	return "check /v1/latest"
}

func (t *latestTest) Run(baseUrl string) error {
	resp, err := makeRequest(baseUrl, "/v1/latest")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("response had non-200 status code %v", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	decodedBody := map[string]interface{}{}

	if err := decoder.Decode(&decodedBody); err != nil {
		return fmt.Errorf("could not decode JSON response: %w", err)
	}

	url, hasUrl := decodedBody["url"]

	if !hasUrl {
		return fmt.Errorf("response body is missing URL: %v", decodedBody)
	}

	if !strings.HasPrefix(url.(string), "https://github.com/batect/batect/releases/tag/") {
		return fmt.Errorf("response body has unexpected value for URL: %s", url)
	}

	return nil
}

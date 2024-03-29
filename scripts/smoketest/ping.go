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
	"fmt"
	"io"
)

type pingTest struct{}

func (t *pingTest) Description() string {
	return "check /ping"
}

func (t *pingTest) Run(baseURL string) error {
	resp, err := makeRequest(baseURL, "/ping")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("response had non-200 status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	bodyText := string(body)

	if bodyText != "pong" {
		return fmt.Errorf("body had unexpected content: %s", bodyText)
	}

	return nil
}

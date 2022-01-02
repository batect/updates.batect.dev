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

package main

import (
	"fmt"
)

type downloadTest struct{}

func (t *downloadTest) Description() string {
	return "check /v1/files/:version/:filename"
}

func (t *downloadTest) Run(baseUrl string) error {
	resp, err := makeRequest(baseUrl, "/v1/files/0.0.0/batect-0.0.0.jar")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return fmt.Errorf("response had non-302 status code %v", resp.StatusCode)
	}

	actualLocation := resp.Header.Get("Location")
	expectedLocation := "https://github.com/batect/batect/releases/download/0.0.0/batect-0.0.0.jar"

	if actualLocation != expectedLocation {
		return fmt.Errorf("response had unexpected location header '%s', expected '%s'", actualLocation, expectedLocation)
	}

	return nil
}

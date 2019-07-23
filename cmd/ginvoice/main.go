/**
 * Copyright (c) 2019-present Future Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"context"
	"errors"
	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"io/ioutil"
	"log"
	"os"
)

const period = 30

// export GOOGLE_APPLICATION_CREDENTIALS=key.json
func main() {
	log.Println("start")

	var (
		projectID     = os.Getenv("GCP_PROJECT")
		tableName     = os.Getenv("TABLE_NAME")
		slackAPIToken = os.Getenv("SLACK_API_TOKEN")
		slackChannel  = os.Getenv("SLACK_CHANNEL")
	)

	if projectID == "" || tableName == "" || slackAPIToken == "" || slackChannel == "" {
		panic(errors.New("missing env"))
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS is required")
	}

	ivc, err := invoice.NewInvoice(projectID)
	if err != nil {
		panic(err)
	}

	costs, err := ivc.FetchBilling(context.Background(), tableName, period)
	if err != nil {
		panic(err)
	}

	plotBytes, err := graph.Draw(costs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("example.png", plotBytes, 0644); err != nil {
		panic(err)
	}

	log.Println("finish")
}

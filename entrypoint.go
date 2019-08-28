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
package gbillingplot

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"github.com/future-architect/gbilling-plot/notify"
	"log"
	"os"
)

const period = 30

type Payload struct {
	Limit int `json:"limit"`
}

func GraphedBilling(ctx context.Context, m *pubsub.Message) error {
	log.Println("start GraphedBilling")

	var (
		projectID    = os.Getenv("GCP_PROJECT")
		tableName    = os.Getenv("TABLE_NAME")
		slackToken   = os.Getenv("SLACK_API_TOKEN")
		slackChannel = os.Getenv("SLACK_CHANNEL")
	)

	if projectID == "" || tableName == "" || slackToken == "" || slackChannel == "" {
		return errors.New("missing env")
	}

	payload, err := decode(m.Data)
	if err != nil {
		log.Printf("error at the fucntion 'decode': %v", err)
		return err
	}
	limit := payload.Limit
	if limit == 0 {
		limit = 8 // default
	}

	ivc, err := invoice.NewInvoice(ctx, projectID)
	if err != nil {
		log.Println("invoice initialize is failed")
		return err
	}

	costs, err := ivc.FetchBilling(ctx, tableName, period)
	if err != nil {
		log.Println("fetch billing is failed")
		return err
	}

	summaryCosts := costs.SummaryLowerProjects(limit)

	plotBytes, err := graph.Draw(summaryCosts)
	if err != nil {
		log.Println("graph draw is failed")
		return err
	}

	notifier := notify.NewSlackNotifier(slackToken, slackChannel)
	if err := notifier.PostImage(ctx, bytes.NewBuffer(plotBytes)); err != nil {
		log.Println("Slack post is failed")
		return err
	}

	log.Println("finish GraphedBilling")
	return nil
}

func decode(payload []byte) (p Payload, err error) {
	if err = json.Unmarshal(payload, &p); err != nil {
		log.Printf("Message[%v] ... Could not decode subscribing data: %v", payload, err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		return
	}
	return
}

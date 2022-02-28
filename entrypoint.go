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
	"context"
	"errors"
	"log"
	"os"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"github.com/future-architect/gbilling-plot/notify"
)

const period = 30
const upperLimit = 5000

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

	plotBytes, err := graph.Draw(costs)
	if err != nil {
		log.Println("graph draw is failed")
		return err
	}

	notifier := notify.NewSlackNotifier(slackToken, slackChannel)

	charts, _ := graph.GetChartValues(costs)
	total := charts[len(charts)-1].Value
	if total > upperLimit {
		log.Println("high cost")
		msg := fmt.Sprintf("<!channel> you must to review current invoice which is higher than %d€: %d", upperLimit, int(total))
		if err := notifier.PostMessage(ctx, msg); err != nil {
			log.Println("Slack post is failed")
			log.Fatal(err)
		}
	}
	if err := notifier.PostImage(ctx, bytes.NewBuffer(plotBytes)); err != nil {
		log.Println("Slack post is failed")
		return err
	}

	log.Println("finish GraphedBilling")
	return nil
}

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
	"errors"
	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"github.com/future-architect/gbilling-plot/notify"
	"log"
	"os"
)

const period = 30

func GraphedBilling(ctx context.Context, msg *pubsub.Message) error {
	log.Println("start GraphedBilling")

	var (
		projectID     = os.Getenv("GCP_PROJECT")
		tableName     = os.Getenv("TABLE_NAME")
		slackAPIToken = os.Getenv("SLACK_API_TOKEN")
		slackChannel  = os.Getenv("SLACK_CHANNEL")
	)

	if projectID == "" || tableName == "" || slackAPIToken == "" || slackChannel == "" {
		return errors.New("missing env")
	}

	ivc, err := invoice.NewInvoice(projectID)
	if err != nil {
		return err
	}

	costs, err := ivc.FetchBilling(ctx, tableName, period)
	if err != nil {
		return err
	}

	plotBytes, err := graph.Draw(costs)
	if err != nil {
		return err
	}

	notifier := notify.NewSlackNotifier(slackAPIToken, slackChannel)
	if err := notifier.PostImage(bytes.NewBuffer(plotBytes)); err != nil {
		return err
	}

	log.Println("finish GraphedBilling")
	return nil
}

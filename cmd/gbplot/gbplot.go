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
	//"bytes"
	"context"
	"flag"
	//"io/ioutil"
	"fmt"
	"log"
	"os"

	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"github.com/future-architect/gbilling-plot/notify"
)

const period = 30

func main() {

	projectID := flag.String("p", os.Getenv("GCP_PROJECT"), "GCP project name")
	tableName := flag.String("t", os.Getenv("TABLE_NAME"), "BigQuery billing table name")
	outFileName := flag.String("o", "out.png", "Output file name")
	sentToSlack := flag.Bool("s", false, "Send to slack")
	flag.StringVar(projectID, "project", "", "GCP project name")
	flag.StringVar(tableName, "table", "", "BigQuery billing table name")
	flag.StringVar(outFileName, "out", "out.png", "Output file name")
	flag.BoolVar(sentToSlack, "slack", false, "Send to slack")
	flag.Parse()

	if *projectID == "" {
		log.Fatal("GCP project name is required")
	}
	if *tableName == "" {
		log.Fatal("BigQuery billing table name is required")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is required")
	}
	slackToken := os.Getenv("SLACK_API_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_API_TOKEN is required")
	}
	slackChannel := os.Getenv("SLACK_CHANNEL")
	if slackChannel == "" {
		log.Fatal("SLACK_CHANNEL is required")
	}

	ctx := context.Background()
	ivc, err := invoice.NewInvoice(ctx, *projectID)
	if err != nil {
		log.Println("invoice initialize is failed")
		log.Fatal(err)
	}

	costs, err := ivc.FetchBilling(ctx, *tableName, period)
	if err != nil {
		log.Println("fetch billing is failed")
		log.Fatal(err)
	}
	charts, _ := graph.GetChartValues(costs)

	notifier := notify.NewSlackNotifier(slackToken, slackChannel)
	total := charts[len(charts)-1].Value
	if total > 50000 {
		log.Println("high cost")
		msg := fmt.Sprintf("<!channel> you must to review current invoice which is higher than %d€: %d", 5000, int(total))
		if err := notifier.PostMessage(ctx, msg); err != nil {
			log.Println("Slack post is failed")
			log.Fatal(err)
		}
	}
	log.Println()
	//plotBytes, err := graph.Draw(costs)
	//if err != nil {
	//	log.Println("graph draw is failed")
	//	log.Fatal(err)
	//}

	//log.Println(plotBytes)

	//notifier := notify.NewSlackNotifier(slackToken, slackChannel)
	//if err := notifier.PostImage(ctx, bytes.NewBuffer(plotBytes)); err != nil {
	//	log.Println("Slack post is failed")
	//	return err
	//}
}

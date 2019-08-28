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
	"flag"
	"github.com/future-architect/gbilling-plot/graph"
	"github.com/future-architect/gbilling-plot/invoice"
	"io/ioutil"
	"log"
	"os"
)

const period = 30

func main() {

	projectID := flag.String("p", os.Getenv("GCP_PROJECT"), "GCP project name")
	tableName := flag.String("t", os.Getenv("TABLE_NAME"), "BigQuery billing table name")
	outFileName := flag.String("o", "out.png", "Output file name")
	flag.StringVar(projectID, "project", "", "GCP project name")
	flag.StringVar(tableName, "table", "", "BigQuery billing table name")
	flag.StringVar(outFileName, "out", "out.png", "Output file name")
	flag.Parse()

	if *projectID == "" || *tableName == "" {
		log.Fatal("missing env")
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is required")
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

	plotBytes, err := graph.Draw(costs)
	if err != nil {
		log.Println("fetch billing is failed")
		log.Fatal(err)
	}


	if err := ioutil.WriteFile(*outFileName, plotBytes, 0644); err != nil {
		log.Fatal(err)
	}

}

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
package invoice

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type invoice struct {
	client *bigquery.Client
}

type Cost struct {
	Project string  `bigquery:"project"`
	Cost    float64 `bigquery:"cost"`
}

type CostList []Cost

func NewInvoice(ctx context.Context, projectID string) (*invoice, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &invoice{client: client}, nil
}

func (i *invoice) FetchBilling(ctx context.Context, tableName string, period int) (CostList, error) {

	endDay := time.Now()
	// beginning of the month
	startDay := time.Date(endDay.Year(), endDay.Month(), 1, 0, 0, 0, 0, endDay.Location())
	stmt := `
		SELECT
			sq.project AS project,
			sq.cost AS cost
		FROM (
			SELECT
				project.id AS project,
				IFNULL(SUM(cost), 0) AS cost
			FROM
				` + "`" + tableName + "`" + `
			WHERE
				DATE(_PARTITIONTIME) BETWEEN ` + startDay.Format("'2006-01-02'") + `
				AND ` + endDay.Format("'2006-01-02'") + `
				AND project.id IS NOT NULL
			GROUP BY project ) AS sq
		ORDER BY
			sq.project
	`
	log.Println(stmt)
	iter, err := i.client.Query(stmt).Read(ctx)
	if err != nil {
		return nil, err
	}

	var costList []Cost
	for {
		var c Cost
		err := iter.Next(&c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		costList = append(costList, c)
	}
	log.Println(costList)
	return costList, nil
}

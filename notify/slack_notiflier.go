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
package notify

import (
	"context"
	"io"

	"github.com/slack-go/slack"
)

type slackNotifier struct {
	slackAPIToken string
	slackChannel  string
}

func NewSlackNotifier(slackAPIToken, slackChannel string) *slackNotifier {
	return &slackNotifier{
		slackAPIToken: slackAPIToken,
		slackChannel:  slackChannel,
	}
}

func (n *slackNotifier) PostImage(ctx context.Context, r io.Reader) error {
	_, err := slack.New(n.slackAPIToken).UploadFileContext(ctx,
		slack.FileUploadParameters{
			Reader:   r,
			Filename: "Stacked Bar Chart on Projects",
			Channels: []string{n.slackChannel},
		})
	return err
}

func (n* slackNotifier) PostMessage(ctx context.Context, msg string) error {
	_, _, err := slack.New(n.slackAPIToken).PostMessageContext(ctx,
		n.slackChannel,
		slack.MsgOptionText(msg, false),
	)
	return err
}

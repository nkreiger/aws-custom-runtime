/*
Copyright 2021 Triggermesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sender

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Sender struct {
	target      string
	contentType string
}

func New(target, contentType string) *Sender {
	return &Sender{
		target:      target,
		contentType: contentType,
	}
}

func (h *Sender) Send(data []byte, statusCode int, writer http.ResponseWriter) error {
	log.Println("noah: begin send flow")

	ctx := context.Background()

	if h.target != "" {
		resp, err := h.request(ctx, data)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to send the data: %w", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadGateway)
			return err
		}

		if body != nil {
			log.Println("response body exists, write back flow")
			return h.reply(ctx, body, statusCode, writer)
		}

		log.Println("respone body is nil, just return status code")
		writer.WriteHeader(statusCode)
		return nil
	}

	return h.reply(ctx, data, statusCode, writer)
}

func (h *Sender) request(ctx context.Context, data []byte) (*http.Response, error) {
	return http.Post(h.target, h.contentType, bytes.NewBuffer(data))
}

func (h *Sender) reply(ctx context.Context, data []byte, statusCode int, writer http.ResponseWriter) error {
	writer.Header().Set("Content-Type", h.contentType)
	log.Println("initial write status code: ", statusCode)
	writer.WriteHeader(statusCode)
	log.Printf("data: %s", data)
	_, err := writer.Write(data)
	log.Println("error is: ", err)
	return err
}

func (h *Sender) Reply (ctx context.Context, data []byte, statusCode int, writer http.ResponseWriter) error {
	return h.reply(ctx, data, statusCode, writer)
}

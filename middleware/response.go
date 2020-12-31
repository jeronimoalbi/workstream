// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"encoding/json"
	"fmt"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
)

func init() {
	middleware.Response(response)
}

// Response handler that resolves the response contents.
func response(r *kusanagi.Response) (*kusanagi.Response, error) {
	httpRequest := r.GetHTTPRequest()
	method := httpRequest.GetMethod()
	if method == "OPTIONS" {
		return r, nil
	}

	var body map[string]interface{}

	httpResponse := r.GetHTTPResponse()
	status := httpResponse.GetStatusCode()
	transport := r.GetTransport()

	// When the request is a C.O.R.S. pre-flight one don't set the body
	if status == 200 && method != "OPTIONS" && method != "DELETE" {
		// Meta information for the response body
		body = map[string]interface{}{
			"meta": map[string]interface{}{
				"id": transport.GetRequestID(),
			},
		}

		// Add the service data to the body
		data := make(map[string]interface{})
		for _, serviceData := range transport.GetData() {
			// Get the action that was called first for the service
			actionData := serviceData.GetActions()[0]
			// Get the data returned by the first call to the service
			data[serviceData.GetName()] = actionData.GetData()[0]
		}
		body["data"] = data
	} else if status >= 400 && status < 500 {
		// Response failed
		body = map[string]interface{}{
			"status": status,
			"error":  string(httpResponse.GetBody()),
		}

		// When one or more action fail add the errors
		errors := make(map[string]interface{})
		for _, e := range transport.GetErrors() {
			errors[e.GetName()] = map[string]interface{}{
				"message": e.GetMessage(),
				"code":    e.GetCode(),
			}
		}

		if len(errors) > 0 {
			body["errors"] = errors
		}
	}

	// When the body is defined add it to the response
	if body != nil {
		// Log the body contents when debug is enabled
		if r.IsDebug() {
			r.Log(body, log.DEBUG)
		}

		// Serialize the body to JSON
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Failed to serialize response body: %v", err)
		}

		// Update the response contents
		httpResponse.SetHeader("Content-Type", "application/json", false)
		httpResponse.SetBody(b)
	}

	return r, nil
}

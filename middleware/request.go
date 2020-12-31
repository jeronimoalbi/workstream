// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"fmt"
	"strconv"
	"strings"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
)

// C.O.R.S. headers for the HTTP response.
var cors = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "GET, POST, PUT, PATCH, DELETE",
	"Access-Control-Allow-Headers": "Content-Type",
}

func init() {
	middleware.Request(request)
}

// Request handler that resolves the URI to service routing.
func request(r *kusanagi.Request) (interface{}, error) {
	// TODO: Read CORS settings from environment (or config file)
	// An OPTIONS request returns the C.O.R.S. headers
	httpRequest := r.GetHTTPRequest()
	method := httpRequest.GetMethod()
	if method == "OPTIONS" {
		response := r.NewResponse(200, "OK")
		httpResponse := response.GetHTTPResponse()
		for name, value := range cors {
			httpResponse.SetHeader(name, value, false)
		}
		return response, nil
	}

	// Split the URL path
	parts := strings.Split(strings.Trim(httpRequest.GetURLPath(), "/"), "/")
	partsCount := len(parts)
	if partsCount < 2 {
		return r.NewResponse(404, "Not Found"), nil
	}

	// URI path must start with "/VERSION/SERVICE".
	// Parse the version (v1, v2, v3, ...).
	version := strings.TrimLeft(parts[0], "v")
	if _, err := strconv.Atoi(version); err != nil {
		r.Log(fmt.Sprintf("Invalid version format in URL path: %s", parts[0]), log.ERROR)
		return r.NewResponse(404, "Not Found"), nil
	}
	// Change the version format to SEMVER
	version = fmt.Sprintf("%s.*.*", version)
	r.SetServiceVersion(version)

	// Set the service name
	service := parts[1]
	r.SetServiceName(service)

	// When an action name is available in the path, use it, otherwise
	// the action name is based on the HTTP request method.
	action := ""
	if partsCount > 3 {
		// The action name MUST be the 4th URI path part, and it is used
		// as the prefix for the action name, followed by the HTTP request
		// method based name.
		// Example URI path for a custom action: /version/service/pk/action
		action = fmt.Sprintf("%s.", parts[3])
	}

	if method == "GET" && partsCount > 2 {
		action += "read"
	} else if method == "GET" {
		action += "list"
	} else if method == "POST" {
		action += "create"
	} else if method == "PUT" {
		action += "replace"
	} else if method == "PATCH" {
		action += "update"
	} else if method == "DELETE" {
		action += "delete"
	} else {
		return r.NewResponse(405, "Method Not Allowed"), nil
	}

	r.SetActionName(action)
	r.Log(fmt.Sprintf(`Resolved service "%s" (%s) action: "%s"`, service, version, action), log.DEBUG)

	return r, nil
}

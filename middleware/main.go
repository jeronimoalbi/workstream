// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"

// KUSANAGI middleware
var middleware = kusanagi.NewMiddleware()

func main() {
	middleware.Run()
}

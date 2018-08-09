// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
)

func environMap() map[string]string {
	m := map[string]string{}

	for _, s := range os.Environ() {
		if n, v := kevArg(s).keyValue(); n != "" {
			m[n] = v
		}
	}

	return m
}

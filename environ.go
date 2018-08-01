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

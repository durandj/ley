package golangci

import (
	"dagger.io/dagger"

	"universe.dagger.io/docker"
	"universe.dagger.io/go"
)

#Lint: {
	// Source code
	source: dagger.#FS

	// golangci-lint version
	version: *"1.46" | string

	// Timeout on how long the linter is allowed to run
	timeout: *"5m" | string

	_image: docker.#Pull & {
		source: "golangci/golangci-lint:v\(version)"
	}

	go.#Container & {
		"source": source
		input: _image.output
		command: {
			name: "golangci-lint"
			flags: {
				run: true
				"--verbose": true
				"--timeout": timeout
				"--color": "never"
			}
		}
	}
}

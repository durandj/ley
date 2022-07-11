package ley

import (
	"dagger.io/dagger"

	"universe.dagger.io/go"

	"github.com/durandj/ley/dagger/golangci"
	"github.com/durandj/ley/dagger/hadolint"
)

dagger.#Plan & {
	client: filesystem: ".": read: {
		contents: dagger.#FS
		include: [
			"go.mod",
			"go.sum",
			".golangci.yml",
			"**/*.go",
			"**/*.sql",
			"Dockerfile",
		]
	}

	actions: {
		_code: client.filesystem.".".read.contents

		lint: {
			go: golangci.#Lint & {
				source: _code
			}

			docker: hadolint.#Lint & {
				source: _code
			}
		}

		test: {
			unit: go.#Test & {
				source: _code
				package: "./..."

				{
					command: flags: "-race": true
				}
			}
		}
	}
}

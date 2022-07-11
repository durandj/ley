package hadolint

import (
	"dagger.io/dagger"

	"universe.dagger.io/docker"
)

#Lint: {
	// Source code
	source: dagger.#FS

	// hadolint version
	version: *"2.9.3" | string

	// Output format for results
	format: *"tty" | ("tty" | "json" | "checkstyle" | "codeclimate" | "gitlab_codeclimate" | "gnu" | "codacy" | "sonarqube" | "sarif")

	// Threshold of errors which can cause a failure
	failure_threshold: *"info" | ("error" | "warning" | "info" | "style" | "ignore" | "none")

	_image: docker.#Pull & {
		source: "hadolint/hadolint:v\(version)"
	}
	_sourcePath: "/src"

	docker.#Run & {
		input: _image.output

		workdir: _sourcePath

		mounts: {
			"source": {
				dest: _sourcePath
				contents: source
			}
		}

		command: {
			name: "hadolint"
			flags: {
				"--verbose": true
				"--no-color": true
				"--format": format
				"--failure-threshold": failure_threshold
				"Dockerfile": true
			}
		}
	}
}

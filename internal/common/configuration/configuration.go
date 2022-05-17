package configuration

// EnvironmentType gives the general category of environment that the
// application is running in. This is helpful for configuring the
// application's behavior based on where it is running.
type EnvironmentType string

const (
	// EnvironmentTypeProduction represents any final stage environment
	// where the application is handling real world traffic.
	EnvironmentTypeProduction EnvironmentType = ""

	// EnvironmentTypeDev represents a development environment where
	// the application is being built.
	EnvironmentTypeDev EnvironmentType = "dev"
)

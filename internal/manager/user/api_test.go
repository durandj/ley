package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	commonconfiguration "github.com/durandj/ley/internal/common/configuration"
	"github.com/durandj/ley/internal/common/rng"
	"github.com/durandj/ley/internal/manager"
	"github.com/durandj/ley/internal/manager/configuration"
	"github.com/durandj/ley/internal/manager/renderable"
	"github.com/durandj/ley/internal/manager/user"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestUserAPIShouldCreateNewUser(t *testing.T) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	createUserRequest := newCreateUserRequest()
	requestBytes, err := json.Marshal(createUserRequest)
	require.Nil(t, err, "should be able to marshal the request body")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/user", serverAddress),
		bytes.NewBuffer(requestBytes),
	)
	require.Nil(t, err, "should be able to create a POST request")
	request.Header.Add("Content-Type", "application/json")

	startTime := time.Now().UTC().Round(time.Second).Add(-2 * time.Second)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusCreated, response.StatusCode)
	endTime := time.Now().UTC().Round(time.Second).Add(2 * time.Second)

	var newUser user.CreateUserResponse
	err = json.NewDecoder(response.Body).Decode(&newUser)
	require.Nil(t, err, "should be able to read response body")

	require.Equal(
		t,
		createUserRequest.Name,
		newUser.Name,
		"should have requested user name",
	)

	createdOn := time.Time(newUser.CreatedOn)
	require.True(
		t,
		startTime.Unix() <= createdOn.Unix() && createdOn.Unix() <= endTime.Unix(),
		"should have a created on timestamp",
	)

	modifiedOn := time.Time(newUser.ModifiedOn)
	require.True(
		t,
		startTime.Unix() <= modifiedOn.Unix() && modifiedOn.Unix() <= endTime.Unix(),
		"should have a modified on timestamp",
	)
}

func TestUserAPIShouldEnforceUsernameUniqueness(t *testing.T) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	createUserRequest := newCreateUserRequest()
	requestBytes, err := json.Marshal(createUserRequest)
	require.Nil(t, err, "should be able to marshal the request body")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/user", serverAddress),
		bytes.NewBuffer(requestBytes),
	)
	require.Nil(t, err, "should be able to create a POST request")
	request.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusCreated, response.StatusCode)

	request, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/user", serverAddress),
		bytes.NewBuffer(requestBytes),
	)
	require.Nil(t, err, "should be able to create a POST request")
	request.Header.Add("Content-Type", "application/json")

	response, err = httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	var creationError renderable.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&creationError)
	require.Nil(t, err, "should be able to read response body")

	require.Equal(
		t,
		"Username is already taken",
		creationError.Message,
		"should have an error message",
	)
}

func TestUserAPIShouldCheckForRequiredFields(t *testing.T) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	testCases := []struct {
		createRequest *user.CreateUserRequest
		errorMessage  string
	}{
		{
			createRequest: newCreateUserRequest().SetName(""),
			errorMessage:  "Unable to create user: Invalid user name ''",
		},
	}

	for _, testCase := range testCases {
		requestBytes, err := json.Marshal(testCase.createRequest)
		require.Nil(t, err, "should be able to marshal the request body")

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("http://%s/user", serverAddress),
			bytes.NewBuffer(requestBytes),
		)
		require.Nil(t, err, "should be able to create a POST request")
		request.Header.Add("Content-Type", "application/json")

		httpClient := http.Client{}
		response, err := httpClient.Do(request)
		require.Nil(t, err, "should be able to complete the request")
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

		var creationError renderable.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&creationError)
		require.Nil(t, err, "should be able to read response body")

		require.Equal(
			t,
			testCase.errorMessage,
			creationError.Message,
			"should have an error message",
		)
	}
}

func TestUserAPIShouldGetAUserByTheirUsername(t *testing.T) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	testUser, err := newTestUser(ctx, serverAddress)
	require.Nil(t, err, "should be able to create a test user")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/user?username=%s", serverAddress, testUser.Name),
		nil,
	)
	require.Nil(t, err, "should be able to create a GET request")
	request.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusOK, response.StatusCode)

	var returnedUser user.GetUserByUsernameResponse
	err = json.NewDecoder(response.Body).Decode(&returnedUser)
	require.Nil(t, err, "should be able to read response body")

	require.Equal(
		t,
		*testUser,
		returnedUser.RenderableUser,
		"should have returned the requested user",
	)
}

func TestUserAPIShouldReturnAnErrorWhenMissingTheUsernameInGetByUsernameRequest(
	t *testing.T,
) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/user", serverAddress),
		nil,
	)
	require.Nil(t, err, "should be able to create a GET request")
	request.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	var getUserError renderable.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&getUserError)
	require.Nil(t, err, "should be able to read response body")

	require.Equal(
		t,
		"Missing query parameter 'username'",
		getUserError.Message,
		"should have an error message",
	)
}

func TestUserAPIShouldReturnAnErrorForANonExistantUsername(
	t *testing.T,
) {
	config, serverAddress := newServiceConfiguration()

	service, err := manager.New(&config)
	require.Nil(t, err, "should be able to create a manager instance")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		_ = service.Run(ctx)
	}()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/user?username=doesnotexist", serverAddress),
		nil,
	)
	require.Nil(t, err, "should be able to create a GET request")
	request.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	require.Nil(t, err, "should be able to complete the request")
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	var getUserError renderable.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&getUserError)
	require.Nil(t, err, "should be able to read response body")

	require.Equal(
		t,
		"Could not find a user with that name",
		getUserError.Message,
		"should have an error message",
	)
}

func newServiceConfiguration() (configuration.Configuration, string) {
	serverHost, serverPort := "localhost", 8080

	config := configuration.Configuration{
		Service: configuration.ServiceConfiguration{
			EnvironmentType: commonconfiguration.EnvironmentTypeDev,
			Host:            serverHost,
			Port:            serverPort,
		},
		Logging: configuration.LoggingConfiguration{
			Level: configuration.LogLevelInfo,
		},
		DB: configuration.DBConfiguration{
			Type: configuration.DBTypePostgres,
			Postgres: configuration.PostgresConfiguration{
				Host:     "127.0.0.1",
				Port:     5432,
				Role:     "ley",
				Password: "ley",
				DBName:   "ley",
				SSLMode:  "disable",
			},
		},
	}

	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)

	return config, serverAddress
}

func newCreateUserRequest() *user.CreateUserRequest {
	surnames := []string{
		"doe",
		"smith",
		"patrick",
		"baird",
		"ng",
		"collins",
		"vance",
		"newton",
		"euler",
		"neumann",
		"edison",
		"tesla",
		"turing",
	}

	initials := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l",
		"m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z",
	}

	firstInitial := rng.RandomChoice(initials)
	middleInitial := rng.RandomChoice(initials)
	lastName := rng.RandomChoice(surnames)

	return &user.CreateUserRequest{
		Name: fmt.Sprintf("%s%s%s", firstInitial, middleInitial, lastName),
	}
}

func newTestUser(ctx context.Context, serverAddress string) (*user.RenderableUser, error) {
	createUserRequest := newCreateUserRequest()
	requestBytes, err := json.Marshal(createUserRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to create test user: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/user", serverAddress),
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to create test user: %w", err)
	}

	request.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Unable to create test user: %w", err)
	}

	if response.StatusCode != http.StatusCreated {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			responseBody = []byte("Unknown error")
		}

		return nil, fmt.Errorf("Unable to create test user: %s", string(responseBody))
	}

	var newUser user.CreateUserResponse
	if err := json.NewDecoder(response.Body).Decode(&newUser); err != nil {
		return nil, fmt.Errorf("Unable to parse response: %w", err)
	}

	var renderedUser user.RenderableUser = newUser.RenderableUser

	return &renderedUser, nil
}

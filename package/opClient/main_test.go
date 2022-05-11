package opClient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/iop/package/system"
	"github.com/vitorqb/iop/package/tempFiles"
	"github.com/vitorqb/iop/package/testUtils"
	"github.com/vitorqb/iop/package/tokenStorage"
)

func TestNewCreatesANewClientInstance(t *testing.T) {
	sys := system.NewMock()
	tokenStorage := tokenStorage.NewInMemoryTokenStorage("")
	opClient := New(&sys, &tokenStorage)
	if opClient.path != DEFAULT_CLIENT {
		t.Fatal("Unexpected path")
	}
}

func TestRunWithTokenAppendsToken(t *testing.T) {
	tokenStorage := tokenStorage.NewInMemoryTokenStorage("FOO")
	opClient := OpClient{
		tokenStorage: &tokenStorage,
		path:         "echo",
	}
	result, err := opClient.runWithToken("bar")
	assert.Nil(t, err)
	assert.Equal(t, string(result), "--session FOO bar\n")
}

func TestEnsureLoggedInSavesTokenUsingTokenStorage(t *testing.T) {
	tempFiles.NewTempScript("#!/bin/sh \necho -n 123").Run(func(scriptPath string) {
		tokenStorage := tokenStorage.NewInMemoryTokenStorage("")
		opClient := OpClient{
			path:         scriptPath,
			tokenStorage: &tokenStorage,
		}
		opClient.EnsureLoggedIn()
		token, _ := opClient.getToken()
		assert.Equal(t, token, "123")
		assert.Equal(t, tokenStorage.Token, "123")
	})
}

func TestEnsureLoggedInExitsIfCmdFails(t *testing.T) {
	mockSystem := system.NewMock()
	tempFiles.NewTempScript("#!/bin/bash \nexit 1").Run(func(scriptPath string) {
		tokenStorage := tokenStorage.NewInMemoryTokenStorage("")
		opClient := OpClient{
			tokenStorage: &tokenStorage,
			sys:          &mockSystem,
			path:         scriptPath,
		}
		opClient.EnsureLoggedIn()
		assert.Equal(t, mockSystem.CrashCallCount, 1)
		assert.Equal(t, mockSystem.LastCrashErrMsg, "Something wen't wrong during signin")
	})
}

func TestGetPasswordRetunsThePassword(t *testing.T) {
	tempFiles.NewTempScript("#!/bin/sh \necho -n '12345\n'").Run(func(scriptPath string) {
		tokenStorage := tokenStorage.NewInMemoryTokenStorage("")
		opClient := OpClient{
			tokenStorage: &tokenStorage,
			path:         scriptPath,
		}
		assert.Equal(t, opClient.GetPassword("itemRef"), "12345")
	})
}

func TestListItemTitlesReturnItemTitles(t *testing.T) {
	testDataFilePath, _ := testUtils.GetTestDataFilePath("op_list_1.json")
	expectedTitles := []string{"some title 1", "some title 2"}
	testFileCatScript := tempFiles.NewTempScript("#!/bin/sh \ncat " + testDataFilePath)
	testFileCatScript.Run(func(scriptPath string) {
		tokenStorage := tokenStorage.NewInMemoryTokenStorage("")
		opClient := OpClient{
			tokenStorage: &tokenStorage,
			path:         scriptPath,
		}
		assert.ElementsMatch(t, expectedTitles, opClient.ListItemTitles())
	})
}

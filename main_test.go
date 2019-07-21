package main

import (
	"os"
	"os/exec"
	"testing"
)

const apigeecli = "./apigeecli"

var token = os.Getenv("APIGEE_TOKEN")
var env = os.Getenv("APIGEE_ENV")
var org = os.Getenv("APIGEE_ORG")

func TestMain(t *testing.T) {
	cmd := exec.Command(apigeecli)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListOrgs(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "list", "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrg(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListEnvs(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEnv(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "get", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

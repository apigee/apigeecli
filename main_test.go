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

// orgs test
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

func TestSetMart(t *testing.T) {
	mart := os.Getenv("MART")
	cmd := exec.Command(apigeecli, "orgs", "setmart", "-o", org, "-m", mart, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetMartWhiteList(t *testing.T) {
	mart := os.Getenv("MART")
	cmd := exec.Command(apigeecli, "orgs", "setmart", "-o", org, "-m", mart, "-w", "false", "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// env tests
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

// developers test
func TestCreateDeveloper(t *testing.T) {
	email := "test@example.com"
	first := "frstname"
	last := "lastname"
	user := "username"

	cmd := exec.Command(apigeecli, "developers", "create", "-o", org, "-n", email, "-f", first, "-s", last, "-u", user, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDeveloper(t *testing.T) {
	email := "test@example.com"

	cmd := exec.Command(apigeecli, "developers", "get", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListDeveloper(t *testing.T) {

	cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListExpandDeveloper(t *testing.T) {
	expand := "true"

	cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-x", expand, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// kvm test
func TestCreateKVM(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateEncKVM(t *testing.T) {
	name := "testEnc"
	enc := "true"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-c", enc, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListKVM(t *testing.T) {

	cmd := exec.Command(apigeecli, "kvms", "list", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteKVM(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEnvKVM(t *testing.T) {
	name := "testEnc"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-e", env, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// sync tests
func TestGetSync(t *testing.T) {

	cmd := exec.Command(apigeecli, "sync", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetSync(t *testing.T) {
	ity := "test@gmail.com"
	cmd := exec.Command(apigeecli, "sync", "set", "-o", org, "-i", ity, "-t", token)
	err := cmd.Run()
	if err == nil {
		t.Fatal("Invalid identity test should have failed")
	}
}

func TestCreateProxy(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "apis", "create", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateProduct(t *testing.T) {
	name := "test"
	envs := "test"
	proxy := "test"
	approval := "auto"

	cmd := exec.Command(apigeecli, "products", "create", "-o", org, "-n", name, "-e", envs, "-p", proxy, "-f", approval, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateApp(t *testing.T) {
	name := "test"
	email := "test@example.com"
	product := "test"

	cmd := exec.Command(apigeecli, "apps", "create", "-o", org, "-n", name, "-e", email, "-p", product, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetApp(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "apps", "get", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestListApp(t *testing.T) {

	cmd := exec.Command(apigeecli, "apps", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteApp(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "apps", "delete", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetProduct(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "products", "list", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListProduct(t *testing.T) {

	cmd := exec.Command(apigeecli, "products", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProduct(t *testing.T) {
	name := "test"

	cmd := exec.Command(apigeecli, "kvms", "create", "-o", org, "-n", name, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeveloper(t *testing.T) {
	email := "test@example.com"

	cmd := exec.Command(apigeecli, "developers", "delete", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

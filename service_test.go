package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var terraformNullResource string = `
resource "null_resource" "example" {
  provisioner "local-exec" {
    command = "echo Hello, World!"
  }
}

`

type testCase struct {
	teste   string
	tasks   []string
	version string
}

func createTemp(location string) (string, error) {
	tempDir, err := os.MkdirTemp("", "testCase*")
	if err != nil {
		return "", err
	}

	temp, err := os.CreateTemp(tempDir, fmt.Sprintf("%s-*.tf", location))
	if err != nil {
		return "", err
	}

	_, err = temp.WriteString(terraformNullResource)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func TestTerraformService(t *testing.T) {
	ctx := context.Background()

	testCase := testCase{
		teste:   "Plan Apply Destroy",
		tasks:   []string{"testCase1", "testCase2", "testCase3"},
		version: "1.9.5",
	}

	t.Run(testCase.teste, func(t *testing.T) {
		svc, err := NewTerraformService(terraformInstaller, testCase.version)
		if err != nil {
			t.Fatalf("terraform service failed: err %v", err)
		}

		defer func() {
			matches, err := filepath.Glob("/tmp/terraform*")
			if err != nil {
				t.Fatal("matches error", err)
			}

			for _, match := range matches {
				os.RemoveAll(match)
			}
		}()

		var locations []string
		for _, taskLocation := range testCase.tasks {
			temp, err := createTemp(taskLocation)
			if err != nil {
				t.Fatalf("create temp failed: %v", err)
			}
			locations = append(locations, temp)
		}

		for _, location := range locations {
			if err := svc.terraformTaskPlan(ctx, location, svc.execPath); err != nil {
				t.Fatalf("terraform service plan failed: err %v", err)
			}

			if err := svc.terraformTaskCreate(ctx, location, svc.execPath); err != nil {
				t.Fatalf("terraform service apply failed: err %v", err)
			}

			find := filepath.Join(location, "terraform.tfstate")
			_, err := os.Stat(find)
			if err != nil {
				if os.IsNotExist(err) {
					t.Fatalf("file does not exist: %v", err)
				}

				t.Fatal()
			}

			if err := svc.terraformTaskDestroy(ctx, location, svc.execPath); err != nil {
				t.Fatalf("terraform service apply failed: err %v", err)
			}

			if err := os.RemoveAll(location); err != nil {
				t.Fatalf("cleanup error: %v", err)
			}
		}
	})
}

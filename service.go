package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const Output = "output"

type TFService struct {
	execPath string
}

type installerFunc func(v string) (string, error)

func NewTerraformService(installer installerFunc, version string) (*TFService, error) {
	slog.Debug("terraform service", "status", "initializing constructor", "version", version)
	execPath, err := installer(version)
	if err != nil {
		return nil, err
	}
	slog.Debug("terraform service", "status", "initialized", "version", version)
	return &TFService{
		execPath: execPath,
	}, nil
}

func terraformInstaller(v string) (string, error) {
	installer := releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion(v)),
	}
	return installer.Install(context.Background())
}

func (t *TFService) terraformTaskPlan(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	_, err = tf.Plan(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *TFService) terraformTaskDestroy(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	return tf.Destroy(ctx)
}

func (t *TFService) terraformTaskCreate(ctx context.Context, wd, execPath string) error {
	tf, err := tfexec.NewTerraform(wd, execPath)
	if err != nil {
		return err
	}

	tf.SetStdout(os.Stdout)
	if err := tf.Init(ctx, tfexec.Reconfigure(true)); err != nil {
		return err
	}

	plan, err := tf.Plan(ctx, tfexec.Out(Output))
	if err != nil {
		return err
	}

	if plan {
		err = tf.Apply(ctx, tfexec.DirOrPlan(Output))
		if err != nil {
			return err
		}
	} else {
		log.Printf("no changes detected in %s", wd)
	}
	return nil
}

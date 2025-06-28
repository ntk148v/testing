package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	execPath, err := installer.Install(ctx)
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := "./learn-terraform-docker-container/"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}
	log.Println("✅ Terraform initialized")
	// Plan
	planOutput, err := tf.Plan(ctx)
	if err != nil {
		log.Fatalf("Terraform plan failed: %v", err)
	}
	if planOutput {
		fmt.Println("✅ Terraform plan completed with changes.")
	} else {
		fmt.Println("✅ Terraform plan completed with no changes.")
	}

	// Apply
	if err := tf.Apply(ctx); err != nil {
		log.Fatalf("Terraform apply failed: %v", err)
	}
	fmt.Println("✅ Terraform apply completed.")

	// Optional: Show outputs
	outputs, err := tf.Output(ctx)
	if err != nil {
		log.Fatalf("Failed to get outputs: %v", err)
	}

	for key, output := range outputs {
		fmt.Printf("%s = %v\n", key, output.Value)
	}

	// Done
	fmt.Println("🚀 Docker container provisioned via terraform-exec.")
}

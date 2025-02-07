package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zinrai/dpbuild/internal/builder"
	"github.com/zinrai/dpbuild/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: dpbuild <command> [<args>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if err := runInit(); err != nil {
			fmt.Printf("Error during initialization: %v\n", err)
			os.Exit(1)
		}
	case "pkg":
		if err := runPackage(); err != nil {
			fmt.Printf("Error during package build: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runInit() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	b := builder.New(cfg)
	if err := b.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	return nil
}

func runPackage() error {
	pkgCmd := flag.NewFlagSet("pkg", flag.ExitOnError)
	dist := pkgCmd.String("dist", "", "Distribution name (e.g., focal, jammy)")
	arch := pkgCmd.String("arch", "", "Architecture (e.g., amd64)")
	dsc := pkgCmd.String("dsc", "", "Path to .dsc file")

	if err := pkgCmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse pkg command arguments: %w", err)
	}

	if *dist == "" || *arch == "" || *dsc == "" {
		return fmt.Errorf("--dist, --arch, and --dsc flags are required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	b := builder.New(cfg)
	if err := b.Package(&builder.PackageOptions{
		Distribution: *dist,
		Architecture: *arch,
		DscFile:      *dsc,
	}); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}

	return nil
}

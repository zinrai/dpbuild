package main

import (
	"flag"
	"fmt"
	"os"
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
	case "source":
		if err := runSource(); err != nil {
			fmt.Printf("Error during source package creation: %v\n", err)
			os.Exit(1)
		}
	case "package":
		if err := runPackage(); err != nil {
			fmt.Printf("Error during package build: %v\n", err)
			os.Exit(1)
		}
	case "update":
		if err := runUpdate(); err != nil {
			fmt.Printf("Error during update: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runInit() error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	b := NewBuilder(cfg)
	if err := b.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	return nil
}

func runSource() error {
	b := &Builder{}
	if err := b.Source(); err != nil {
		return fmt.Errorf("failed to create source package: %w", err)
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

	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	b := NewBuilder(cfg)
	if err := b.Package(&PackageOptions{
		Distribution: *dist,
		Architecture: *arch,
		DscFile:      *dsc,
	}); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}

	return nil
}

func runUpdate() error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	b := NewBuilder(cfg)
	if err := b.Update(); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	return nil
}

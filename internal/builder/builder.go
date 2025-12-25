package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zinrai/dpbuild/internal/config"
)

type Builder struct {
	config *config.Config
}

func New(cfg *config.Config) *Builder {
	return &Builder{
		config: cfg,
	}
}

func (b *Builder) pbuilderDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".dpbuild", "pbuilder"), nil
}

func (b *Builder) Init() error {
	pbuilderDir, err := b.pbuilderDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(pbuilderDir, 0755); err != nil {
		return fmt.Errorf("failed to create pbuilder directory: %w", err)
	}

	if err := b.checkPbuilder(); err != nil {
		return err
	}

	for _, env := range b.config.Environments {
		fmt.Printf("Setting up environment for %s (%s)...\n", env.Distribution, env.Architecture)

		baseFileName := fmt.Sprintf("%s-%s.tgz", env.Distribution, env.Architecture)
		basePath := filepath.Join(pbuilderDir, baseFileName)

		if _, err := os.Stat(basePath); err == nil {
			fmt.Printf("Base image already exists for %s (%s), skipping...\n", env.Distribution, env.Architecture)
			continue
		}

		args := []string{
			"pbuilder", "create",
			"--distribution", env.Distribution,
			"--architecture", env.Architecture,
			"--mirror", env.Mirror,
			"--components", strings.Join(env.Components, " "),
			"--basetgz", basePath,
		}

		fmt.Printf("Executing: sudo %s\n", strings.Join(args, " "))

		cmd := exec.Command("sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run pbuilder create for %s: %w", env.Distribution, err)
		}

		fmt.Printf("Successfully initialized environment for %s (%s)\n", env.Distribution, env.Architecture)
	}

	return nil
}

func (b *Builder) checkPbuilder() error {
	_, err := exec.LookPath("/usr/sbin/pbuilder")
	if err != nil {
		return fmt.Errorf("pbuilder command not found: %w", err)
	}
	return nil
}

type PackageOptions struct {
	Distribution string
	Architecture string
	DscFile      string
}

func (b *Builder) Package(opts *PackageOptions) error {
	if err := b.validatePackageOptions(opts); err != nil {
		return err
	}

	if err := b.checkPbuilder(); err != nil {
		return err
	}

	if _, err := os.Stat("debian"); err != nil {
		return fmt.Errorf("debian directory not found in current directory")
	}

	if _, err := os.Stat(opts.DscFile); err != nil {
		return fmt.Errorf("dsc file not found: %s", opts.DscFile)
	}

	outputDir := filepath.Join(".packages", fmt.Sprintf("%s-%s", opts.Distribution, opts.Architecture))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	pbuilderDir, err := b.pbuilderDir()
	if err != nil {
		return err
	}

	baseFileName := fmt.Sprintf("%s-%s.tgz", opts.Distribution, opts.Architecture)
	basePath := filepath.Join(pbuilderDir, baseFileName)

	args := []string{
		"pbuilder", "build",
		"--basetgz", basePath,
		"--buildresult", outputDir,
		opts.DscFile,
	}

	fmt.Printf("Executing: sudo %s\n", strings.Join(args, " "))

	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "."

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pbuilder build: %w", err)
	}

	return nil
}

func (b *Builder) validatePackageOptions(opts *PackageOptions) error {
	valid := false
	for _, env := range b.config.Environments {
		if env.Distribution == opts.Distribution && env.Architecture == opts.Architecture {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid distribution and architecture combination: %s-%s", opts.Distribution, opts.Architecture)
	}
	return nil
}

func (b *Builder) Update() error {
	if err := b.checkPbuilder(); err != nil {
		return err
	}

	pbuilderDir, err := b.pbuilderDir()
	if err != nil {
		return err
	}

	for _, env := range b.config.Environments {
		baseFileName := fmt.Sprintf("%s-%s.tgz", env.Distribution, env.Architecture)
		basePath := filepath.Join(pbuilderDir, baseFileName)

		if _, err := os.Stat(basePath); err != nil {
			fmt.Printf("Base image not found for %s (%s), skipping...\n", env.Distribution, env.Architecture)
			continue
		}

		fmt.Printf("Updating environment for %s (%s)...\n", env.Distribution, env.Architecture)

		args := []string{
			"pbuilder", "update",
			"--basetgz", basePath,
		}

		fmt.Printf("Executing: sudo %s\n", strings.Join(args, " "))

		cmd := exec.Command("sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run pbuilder update for %s: %w", env.Distribution, err)
		}

		fmt.Printf("Successfully updated environment for %s (%s)\n", env.Distribution, env.Architecture)
	}

	return nil
}

func (b *Builder) Source() error {
	if _, err := os.Stat("debian"); err != nil {
		return fmt.Errorf("debian directory not found in current directory")
	}

	versionCmd := exec.Command("dpkg-parsechangelog", "-S", "Version")
	versionOut, err := versionCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get package version: %w", err)
	}
	fullVersion := strings.TrimSpace(string(versionOut))

	version := strings.Split(fullVersion, "-")[0]

	sourceCmd := exec.Command("dpkg-parsechangelog", "-S", "Source")
	sourceOut, err := sourceCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get package source name: %w", err)
	}
	sourceName := strings.TrimSpace(string(sourceOut))

	origFileName := fmt.Sprintf("../%s_%s.orig.tar.xz", sourceName, version)

	args := []string{
		"tar", "cJf", origFileName,
		"--exclude", "./debian",
		"--exclude", "./.git",
		"--transform", fmt.Sprintf("s,^\\./,%s-%s/,", sourceName, version),
		".",
	}

	fmt.Printf("Executing: %s\n", strings.Join(args, " "))
	tarCmd := exec.Command(args[0], args[1:]...)
	tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr
	if err := tarCmd.Run(); err != nil {
		return fmt.Errorf("failed to create orig.tar.xz: %w", err)
	}

	sourceArgs := []string{"dpkg-source", "-b", "."}
	fmt.Printf("Executing: %s\n", strings.Join(sourceArgs, " "))
	sourceCmd = exec.Command(sourceArgs[0], sourceArgs[1:]...)
	sourceCmd.Stdout = os.Stdout
	sourceCmd.Stderr = os.Stderr
	if err := sourceCmd.Run(); err != nil {
		return fmt.Errorf("failed to run dpkg-source: %w", err)
	}

	return nil
}

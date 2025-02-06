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

func (b *Builder) Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pbuilderDir := filepath.Join(home, "pbuilder")
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

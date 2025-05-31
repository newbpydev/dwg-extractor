package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCrossPlatformBuild tests cross-platform build capabilities
func TestCrossPlatformBuild(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
	}{
		{
			name:     "Windows AMD64 build",
			goos:     "windows",
			goarch:   "amd64",
			expected: "go-dwg-extractor.exe",
		},
		{
			name:     "Linux AMD64 build",
			goos:     "linux",
			goarch:   "amd64",
			expected: "go-dwg-extractor",
		},
		{
			name:     "macOS AMD64 build",
			goos:     "darwin",
			goarch:   "amd64",
			expected: "go-dwg-extractor",
		},
		{
			name:     "macOS ARM64 build",
			goos:     "darwin",
			goarch:   "arm64",
			expected: "go-dwg-extractor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement BuildManager
			buildManager := NewBuildManager()
			require.NotNil(t, buildManager, "Expected build manager to be created")

			// Create build configuration
			config := BuildConfig{
				GOOS:       tt.goos,
				GOARCH:     tt.goarch,
				OutputName: tt.expected,
				SourcePath: ".",
				OutputDir:  "dist",
				Version:    "1.0.0",
				BuildTime:  "2024-12-28T12:00:00Z",
				GitCommit:  "abcd1234",
			}

			// Execute build
			result, err := buildManager.Build(config)
			require.NoError(t, err, "Expected build to succeed")

			// Verify build result
			assert.Equal(t, tt.expected, result.OutputFile, "Expected correct output filename")
			assert.Equal(t, tt.goos, result.GOOS, "Expected correct GOOS")
			assert.Equal(t, tt.goarch, result.GOARCH, "Expected correct GOARCH")
			assert.True(t, result.Success, "Expected build to be successful")
			assert.Greater(t, result.BuildTime, int64(0), "Expected positive build time")
		})
	}
}

// TestBuildScript tests build script generation and execution
func TestBuildScript(t *testing.T) {
	tests := []struct {
		name             string
		scriptType       ScriptType
		platforms        []Platform
		expectedCommands []string
		expectedFiles    []string
	}{
		{
			name:       "PowerShell build script for Windows",
			scriptType: ScriptTypePowerShell,
			platforms: []Platform{
				{GOOS: "windows", GOARCH: "amd64"},
				{GOOS: "linux", GOARCH: "amd64"},
			},
			expectedCommands: []string{
				"$env:GOOS='windows'",
				"$env:GOARCH='amd64'",
				"go build",
			},
			expectedFiles: []string{"build.ps1"},
		},
		{
			name:       "Bash build script for Unix",
			scriptType: ScriptTypeBash,
			platforms: []Platform{
				{GOOS: "linux", GOARCH: "amd64"},
				{GOOS: "darwin", GOARCH: "amd64"},
			},
			expectedCommands: []string{
				"export GOOS=linux",
				"export GOARCH=amd64",
				"go build",
			},
			expectedFiles: []string{"build.sh"},
		},
		{
			name:       "Makefile build script",
			scriptType: ScriptTypeMakefile,
			platforms: []Platform{
				{GOOS: "windows", GOARCH: "amd64"},
				{GOOS: "linux", GOARCH: "amd64"},
				{GOOS: "darwin", GOARCH: "amd64"},
			},
			expectedCommands: []string{
				"build-windows:",
				"build-linux:",
				"build-darwin:",
				"GOOS=",
				"GOARCH=",
			},
			expectedFiles: []string{"Makefile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement ScriptGenerator
			scriptGenerator := NewScriptGenerator()
			require.NotNil(t, scriptGenerator, "Expected script generator to be created")

			// Generate build script
			script, err := scriptGenerator.GenerateBuildScript(tt.scriptType, tt.platforms)
			require.NoError(t, err, "Expected script generation to succeed")

			// Verify script content
			assert.NotEmpty(t, script.Content, "Expected script content to be generated")
			assert.Equal(t, tt.scriptType, script.Type, "Expected correct script type")
			assert.Equal(t, tt.expectedFiles[0], script.Filename, "Expected correct filename")

			// Check for expected commands
			for _, expectedCmd := range tt.expectedCommands {
				assert.Contains(t, script.Content, expectedCmd,
					"Expected script to contain command: %s", expectedCmd)
			}

			// Verify script is executable
			assert.True(t, script.IsExecutable, "Expected script to be marked as executable")
		})
	}
}

// TestVersioning tests version injection during build
func TestVersioning(t *testing.T) {
	tests := []struct {
		name            string
		version         string
		gitCommit       string
		buildTime       string
		expectedLdflags []string
	}{
		{
			name:      "Semantic version with git commit",
			version:   "v1.2.3",
			gitCommit: "abc123def456",
			buildTime: "2024-12-28T12:00:00Z",
			expectedLdflags: []string{
				"-X main.version=v1.2.3",
				"-X main.gitCommit=abc123def456",
				"-X main.buildTime=2024-12-28T12:00:00Z",
			},
		},
		{
			name:      "Development version",
			version:   "v0.0.0-dev",
			gitCommit: "dirty",
			buildTime: "2024-12-28T12:00:00Z",
			expectedLdflags: []string{
				"-X main.version=v0.0.0-dev",
				"-X main.gitCommit=dirty",
				"-X main.buildTime=2024-12-28T12:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement VersionManager
			versionManager := NewVersionManager()
			require.NotNil(t, versionManager, "Expected version manager to be created")

			// Create version info
			versionInfo := VersionInfo{
				Version:   tt.version,
				GitCommit: tt.gitCommit,
				BuildTime: tt.buildTime,
			}

			// Generate ldflags
			ldflags := versionManager.GenerateLdflags(versionInfo)
			assert.NotEmpty(t, ldflags, "Expected ldflags to be generated")

			// Check for expected ldflags
			for _, expectedFlag := range tt.expectedLdflags {
				assert.Contains(t, ldflags, expectedFlag,
					"Expected ldflags to contain: %s", expectedFlag)
			}
		})
	}
}

// TestPackaging tests creation of distribution packages
func TestPackaging(t *testing.T) {
	tests := []struct {
		name           string
		platform       Platform
		files          []string
		packageType    PackageType
		expectedOutput string
		expectedFiles  []string
	}{
		{
			name:     "Windows ZIP package",
			platform: Platform{GOOS: "windows", GOARCH: "amd64"},
			files: []string{
				"go-dwg-extractor.exe",
				"README.md",
				"LICENSE",
			},
			packageType:    PackageTypeZip,
			expectedOutput: "go-dwg-extractor-windows-amd64.zip",
			expectedFiles: []string{
				"go-dwg-extractor.exe",
				"README.md",
				"LICENSE",
			},
		},
		{
			name:     "Linux tarball package",
			platform: Platform{GOOS: "linux", GOARCH: "amd64"},
			files: []string{
				"go-dwg-extractor",
				"README.md",
				"LICENSE",
			},
			packageType:    PackageTypeTarGz,
			expectedOutput: "go-dwg-extractor-linux-amd64.tar.gz",
			expectedFiles: []string{
				"go-dwg-extractor",
				"README.md",
				"LICENSE",
			},
		},
		{
			name:     "macOS ZIP package",
			platform: Platform{GOOS: "darwin", GOARCH: "amd64"},
			files: []string{
				"go-dwg-extractor",
				"README.md",
				"LICENSE",
			},
			packageType:    PackageTypeZip,
			expectedOutput: "go-dwg-extractor-darwin-amd64.zip",
			expectedFiles: []string{
				"go-dwg-extractor",
				"README.md",
				"LICENSE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement PackageManager
			packageManager := NewPackageManager()
			require.NotNil(t, packageManager, "Expected package manager to be created")

			// Create package configuration
			config := PackageConfig{
				Platform:    tt.platform,
				Files:       tt.files,
				PackageType: tt.packageType,
				OutputDir:   "dist",
				Version:     "v1.0.0",
			}

			// Create package
			packageResult, err := packageManager.CreatePackage(config)
			require.NoError(t, err, "Expected package creation to succeed")

			// Verify package result
			assert.Equal(t, tt.expectedOutput, packageResult.Filename, "Expected correct package filename")
			assert.Equal(t, tt.packageType, packageResult.Type, "Expected correct package type")
			assert.ElementsMatch(t, tt.expectedFiles, packageResult.IncludedFiles, "Expected correct files in package")
			assert.Greater(t, packageResult.Size, int64(0), "Expected package to have positive size")
		})
	}
}

// TestBuildPipeline tests the complete build pipeline
func TestBuildPipeline(t *testing.T) {
	tests := []struct {
		name              string
		config            PipelineConfig
		expectedBuilds    int
		expectedPackages  int
		expectedArtifacts []string
	}{
		{
			name: "Multi-platform release pipeline",
			config: PipelineConfig{
				Version: "v1.0.0",
				Platforms: []Platform{
					{GOOS: "windows", GOARCH: "amd64"},
					{GOOS: "linux", GOARCH: "amd64"},
					{GOOS: "darwin", GOARCH: "amd64"},
				},
				IncludeFiles:   []string{"README.md", "LICENSE"},
				OutputDir:      "dist",
				CreatePackages: true,
			},
			expectedBuilds:   3,
			expectedPackages: 3,
			expectedArtifacts: []string{
				"go-dwg-extractor-windows-amd64.zip",
				"go-dwg-extractor-linux-amd64.tar.gz",
				"go-dwg-extractor-darwin-amd64.zip",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement BuildPipeline
			pipeline := NewBuildPipeline()
			require.NotNil(t, pipeline, "Expected build pipeline to be created")

			// Execute pipeline
			result, err := pipeline.Execute(tt.config)
			require.NoError(t, err, "Expected pipeline execution to succeed")

			// Verify pipeline result
			assert.Equal(t, tt.expectedBuilds, len(result.Builds), "Expected correct number of builds")
			assert.Equal(t, tt.expectedPackages, len(result.Packages), "Expected correct number of packages")
			assert.True(t, result.Success, "Expected pipeline to succeed")

			// Check for expected artifacts
			for _, expectedArtifact := range tt.expectedArtifacts {
				found := false
				for _, pkg := range result.Packages {
					if pkg.Filename == expectedArtifact {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected artifact %s to be created", expectedArtifact)
			}
		})
	}
}

// TestBuildEnvironment tests build environment validation
func TestBuildEnvironment(t *testing.T) {
	tests := []struct {
		name            string
		requiredTools   []string
		expectedValid   bool
		expectedMissing []string
	}{
		{
			name:            "Valid Go environment",
			requiredTools:   []string{"go", "git"},
			expectedValid:   true,
			expectedMissing: []string{},
		},
		{
			name:            "Missing tools",
			requiredTools:   []string{"go", "git", "nonexistent-tool"},
			expectedValid:   false,
			expectedMissing: []string{"nonexistent-tool"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement EnvironmentValidator
			validator := NewEnvironmentValidator()
			require.NotNil(t, validator, "Expected environment validator to be created")

			// Validate environment
			result := validator.ValidateEnvironment(tt.requiredTools)

			assert.Equal(t, tt.expectedValid, result.IsValid, "Expected correct validation result")
			assert.ElementsMatch(t, tt.expectedMissing, result.MissingTools, "Expected correct missing tools")

			if result.IsValid {
				assert.NotEmpty(t, result.GoVersion, "Expected Go version to be detected")
				assert.NotEmpty(t, result.Environment, "Expected environment info to be populated")
			}
		})
	}
}

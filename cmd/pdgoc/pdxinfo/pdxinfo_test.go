package pdxinfo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/playdate-go/pdgo/cmd/pdgoc/config"
)

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		check   func(t *testing.T, content string)
	}{
		{
			name: "minimal config",
			cfg: &config.Config{
				Meta: &config.Meta{
					Name:        "TestGame",
					Author:      "TestAuthor",
					Desc:        "A test game",
					BundleID:    "com.test.game",
					Version:     "1.0.0",
					BuildNumber: "1",
				},
			},
			wantErr: false,
			check: func(t *testing.T, content string) {
				required := []string{
					"name=TestGame",
					"author=TestAuthor",
					"description=A test game",
					"bundleID=com.test.game",
					"version=1.0.0",
					"buildNumber=1",
				}
				for _, r := range required {
					if !strings.Contains(content, r) {
						t.Errorf("pdxinfo missing required field: %s", r)
					}
				}
			},
		},
		{
			name: "full config with optional fields",
			cfg: &config.Config{
				Meta: &config.Meta{
					Name:            "FullGame",
					Author:          "FullAuthor",
					Desc:            "A full game",
					BundleID:        "com.full.game",
					Version:         "2.0.0",
					BuildNumber:     "42",
					ImagePath:       "icon.png",
					LaunchSoundPath: "launch.wav",
					ContentWarn:     "Warning 1",
					ContentWarn2:    "Warning 2",
				},
			},
			wantErr: false,
			check: func(t *testing.T, content string) {
				optional := []string{
					"imagePath=icon.png",
					"launchSoundPath=launch.wav",
					"contentWarning=Warning 1",
					"contentWarning2=Warning 2",
				}
				for _, o := range optional {
					if !strings.Contains(content, o) {
						t.Errorf("pdxinfo missing optional field: %s", o)
					}
				}
			},
		},
		{
			name: "config with empty optional fields",
			cfg: &config.Config{
				Meta: &config.Meta{
					Name:            "MinimalGame",
					Author:          "Author",
					Desc:            "Desc",
					BundleID:        "com.minimal",
					Version:         "1.0",
					BuildNumber:     "1",
					ImagePath:       "",
					LaunchSoundPath: "",
					ContentWarn:     "",
					ContentWarn2:    "",
				},
			},
			wantErr: false,
			check: func(t *testing.T, content string) {
				if strings.Contains(content, "imagePath=") {
					t.Error("pdxinfo should not contain empty imagePath")
				}
				if strings.Contains(content, "launchSoundPath=") {
					t.Error("pdxinfo should not contain empty launchSoundPath")
				}
				if strings.Contains(content, "contentWarning=") {
					t.Error("pdxinfo should not contain empty contentWarning")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "pdxinfo-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			err = CreateFile(tt.cfg, tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				pdxInfoPath := filepath.Join(tmpDir, "pdxinfo")
				content, err := os.ReadFile(pdxInfoPath)
				if err != nil {
					t.Fatalf("failed to read pdxinfo file: %v", err)
				}

				if tt.check != nil {
					tt.check(t, string(content))
				}
			}
		})
	}
}

func TestCreateFileNilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("CreateFile with nil config should panic")
		}
	}()

	tmpDir, _ := os.MkdirTemp("", "pdxinfo-test-*")
	defer os.RemoveAll(tmpDir)

	CreateFile(nil, tmpDir)
}

func TestCreateFileInvalidPath(t *testing.T) {
	cfg := &config.Config{
		Meta: &config.Meta{
			Name:        "Test",
			Author:      "Test",
			Desc:        "Test",
			BundleID:    "com.test",
			Version:     "1.0",
			BuildNumber: "1",
		},
	}

	err := CreateFile(cfg, "/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("CreateFile with invalid path should return error")
	}
}

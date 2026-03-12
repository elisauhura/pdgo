package config

import (
	"testing"
)

func TestMetaValidate(t *testing.T) {
	tests := []struct {
		name    string
		meta    *Meta
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil meta",
			meta:    nil,
			wantErr: true,
			errMsg:  "passed nil config.Meta struct",
		},
		{
			name: "missing name",
			meta: &Meta{
				Author:      "test",
				Desc:        "test",
				BundleID:    "com.test",
				Version:     "1.0",
				BuildNumber: "1",
			},
			wantErr: true,
			errMsg:  "'name' pdxinfo property is mandatory",
		},
		{
			name: "missing author",
			meta: &Meta{
				Name:        "test",
				Desc:        "test",
				BundleID:    "com.test",
				Version:     "1.0",
				BuildNumber: "1",
			},
			wantErr: true,
			errMsg:  "'author' pdxinfo property is mandatory",
		},
		{
			name: "missing desc",
			meta: &Meta{
				Name:        "test",
				Author:      "test",
				BundleID:    "com.test",
				Version:     "1.0",
				BuildNumber: "1",
			},
			wantErr: true,
			errMsg:  "'description' pdxinfo property is mandatory",
		},
		{
			name: "missing bundleID",
			meta: &Meta{
				Name:        "test",
				Author:      "test",
				Desc:        "test",
				Version:     "1.0",
				BuildNumber: "1",
			},
			wantErr: true,
			errMsg:  "'bundleID' pdxinfo property is mandatory",
		},
		{
			name: "missing version",
			meta: &Meta{
				Name:        "test",
				Author:      "test",
				Desc:        "test",
				BundleID:    "com.test",
				BuildNumber: "1",
			},
			wantErr: true,
			errMsg:  "'version' pdxinfo property is mandatory",
		},
		{
			name: "missing buildNumber",
			meta: &Meta{
				Name:     "test",
				Author:   "test",
				Desc:     "test",
				BundleID: "com.test",
				Version:  "1.0",
			},
			wantErr: true,
			errMsg:  "'buildNumber' pdxinfo property is mandatory",
		},
		{
			name: "valid minimal config",
			meta: &Meta{
				Name:        "test",
				Author:      "test",
				Desc:        "test",
				BundleID:    "com.test",
				Version:     "1.0",
				BuildNumber: "1",
			},
			wantErr: false,
		},
		{
			name: "valid full config",
			meta: &Meta{
				Name:            "test",
				Author:          "test",
				Desc:            "test",
				BundleID:        "com.test",
				Version:         "1.0",
				BuildNumber:     "1",
				ImagePath:       "img.png",
				LaunchSoundPath: "sound.wav",
				ContentWarn:     "warning",
				ContentWarn2:    "warning2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.meta.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Meta.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Meta.Validate() error = %v, want containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestSystemValidate(t *testing.T) {
	tests := []struct {
		name    string
		system  *System
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no mode specified",
			system:  &System{},
			wantErr: true,
			errMsg:  "at least '-sim' or '-device' or '-run' must be defined",
		},
		{
			name: "deploy without device (only deploy set)",
			system: &System{
				DeployMode: true,
			},
			wantErr: true,
			errMsg:  "at least '-sim' or '-device' or '-run' must be defined",
		},
		{
			name: "deploy with sim but without device",
			system: &System{
				SimMode:    true,
				DeployMode: true,
			},
			wantErr: true,
			errMsg:  "'-deploy' requires '-device' flag",
		},
		{
			name: "sim mode only",
			system: &System{
				SimMode: true,
			},
			wantErr: false,
		},
		{
			name: "device mode only",
			system: &System{
				DeviceMode: true,
			},
			wantErr: false,
		},
		{
			name: "run mode only",
			system: &System{
				RunMode: true,
			},
			wantErr: false,
		},
		{
			name: "deploy with device",
			system: &System{
				DeviceMode: true,
				DeployMode: true,
			},
			wantErr: false,
		},
		{
			name: "sim and device",
			system: &System{
				SimMode:    true,
				DeviceMode: true,
			},
			wantErr: false,
		},
		{
			name: "all modes",
			system: &System{
				SimMode:    true,
				DeviceMode: true,
				RunMode:    true,
				DeployMode: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.system.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("System.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("System.Validate() error = %v, want containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestMetaString(t *testing.T) {
	meta := &Meta{
		Name:        "test",
		Author:      "author",
		Desc:        "desc",
		BundleID:    "com.test",
		Version:     "1.0",
		BuildNumber: "1",
	}

	result := meta.String()
	if result == "" {
		t.Error("Meta.String() returned empty string")
	}
	if !contains(result, "test") {
		t.Errorf("Meta.String() = %v, want containing 'test'", result)
	}
}

func TestSystemString(t *testing.T) {
	system := &System{
		SimMode:    true,
		DeviceMode: false,
		RunMode:    true,
	}

	result := system.String()
	if result == "" {
		t.Error("System.String() returned empty string")
	}
}

func TestConfigString(t *testing.T) {
	cfg := &Config{
		System: &System{SimMode: true},
		Meta:   &Meta{Name: "test"},
	}

	result := cfg.String()
	if result == "" {
		t.Error("Config.String() returned empty string")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package core

import "testing"

func TestConfigKeysInfo_AllKeysHaveDescriptions(t *testing.T) {
	for key, entry := range configKeyRegistry {
		if entry.Description == "" {
			t.Errorf("config key %q has no description", key)
		}
	}
}

func TestConfigKeysInfo_UnsetValueShowsEmpty(t *testing.T) {
	v := setupTestVault(t)
	infos := v.ConfigKeysInfo()
	for _, info := range infos {
		if info.Key == "cli.default_type" {
			if info.Value != "" {
				t.Errorf("expected unset cli.default_type to have empty value, got %q", info.Value)
			}
			return
		}
	}
	t.Fatal("cli.default_type not found in ConfigKeysInfo")
}

func TestConfigKeysInfo_SetValueShowsCurrent(t *testing.T) {
	v := setupTestVault(t)
	if err := v.SetConfigValue("cli.default_type", "idea"); err != nil {
		t.Fatal(err)
	}
	infos := v.ConfigKeysInfo()
	for _, info := range infos {
		if info.Key == "cli.default_type" {
			if info.Value != "idea" {
				t.Errorf("expected cli.default_type value to be %q, got %q", "idea", info.Value)
			}
			if info.Default != "" {
				t.Errorf("expected cli.default_type default to be empty, got %q", info.Default)
			}
			return
		}
	}
	t.Fatal("cli.default_type not found in ConfigKeysInfo")
}

func TestConfigKeysInfo_ReturnsSorted(t *testing.T) {
	v := setupTestVault(t)
	infos := v.ConfigKeysInfo()
	for i := 1; i < len(infos); i++ {
		if infos[i].Key < infos[i-1].Key {
			t.Errorf("keys not sorted: %q comes after %q", infos[i].Key, infos[i-1].Key)
		}
	}
}

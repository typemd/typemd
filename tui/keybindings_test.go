package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func TestValidateKeyString_Valid(t *testing.T) {
	cases := []string{
		"a", "/", "=", "?", ",", ".",
		"ctrl+s", "ctrl+e", "ctrl+c",
		"ctrl+shift+a", "shift+ctrl+a", // modifier order accepted either way
		"alt+x",
		"esc", "enter", "tab", "space",
		"up", "down", "left", "right",
		"pgup", "pgdown", "home", "end",
		"f1", "f12",
	}
	for _, c := range cases {
		if err := validateKeyString(c); err != nil {
			t.Errorf("validateKeyString(%q) unexpected error: %v", c, err)
		}
	}
}

func TestValidateKeyString_Invalid(t *testing.T) {
	cases := []struct {
		input string
		// substring that must appear in the error message
		reason string
	}{
		{"", "empty"},
		{"crtl+s", "unknown"},       // typo
		{"ctrl+", "empty segment"},  // trailing +
		{"+a", "empty segment"},     // leading +
		{"ctrl++s", "empty segment"}, // double +
		{"ctrl+ctrl+a", "duplicate"}, // duplicate modifier
		{"notakey", "unknown key"},
		{"ctrl+notakey", "unknown key"},
	}
	for _, c := range cases {
		err := validateKeyString(c.input)
		if err == nil {
			t.Errorf("validateKeyString(%q) expected error, got nil", c.input)
			continue
		}
		if !strings.Contains(err.Error(), c.reason) {
			t.Errorf("validateKeyString(%q) error = %q, want containing %q", c.input, err.Error(), c.reason)
		}
	}
}

func TestBuildKeyMap_NoConfigReturnsDefaults(t *testing.T) {
	km, issues := buildKeyMap(nil)
	if len(issues) != 0 {
		t.Errorf("buildKeyMap(nil) issues = %v, want empty", issues)
	}
	// Spot-check a few defaults.
	if got := km.Stats.Help().Key; got != "ctrl+s" {
		t.Errorf("Stats help key = %q, want %q", got, "ctrl+s")
	}
	if got := km.Settings.Help().Key; got != "," {
		t.Errorf("Settings help key = %q, want %q", got, ",")
	}
}

func TestBuildKeyMap_OverrideOneAction(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionStats: "ctrl+d",
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 0 {
		t.Fatalf("buildKeyMap override issues = %v, want empty", issues)
	}
	if got := km.Stats.Help().Key; got != "ctrl+d" {
		t.Errorf("Stats help key after override = %q, want %q", got, "ctrl+d")
	}
	// translate("ctrl+d") should route to default "ctrl+s" so dispatch still works.
	if got := km.translate("ctrl+d"); got != "ctrl+s" {
		t.Errorf("translate(ctrl+d) = %q, want %q (should route to default)", got, "ctrl+s")
	}
	// The original default key "ctrl+s" was rebound away; translate must
	// return the unbound sentinel so the dispatch switch no longer matches
	// `case "ctrl+s"` and stats does not open (issue #396).
	if got := km.translate("ctrl+s"); got != unboundSentinel {
		t.Errorf("translate(ctrl+s) after override = %q, want unboundSentinel", got)
	}
	// And the resolved binding's keys list no longer contains the old default.
	for _, k := range km.Stats.Keys() {
		if k == "ctrl+s" {
			t.Errorf("after override, Stats binding keys should not contain ctrl+s, got %v", km.Stats.Keys())
		}
	}
}

func TestBuildKeyMap_EmptyStringMeansDefault(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionSearch: "",
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 0 {
		t.Errorf("buildKeyMap empty-string issues = %v, want empty", issues)
	}
	if got := km.Search.Help().Key; got != "/" {
		t.Errorf("Search help key with empty override = %q, want default %q", got, "/")
	}
}

func TestBuildKeyMap_UnknownAction(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		"not_a_real_action": "ctrl+x",
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 1 {
		t.Fatalf("issues len = %d, want 1. issues=%v", len(issues), issues)
	}
	if issues[0].Kind != IssueUnknownAction {
		t.Errorf("issue[0].Kind = %v, want IssueUnknownAction", issues[0].Kind)
	}
	if issues[0].Action != "not_a_real_action" {
		t.Errorf("issue[0].Action = %q, want %q", issues[0].Action, "not_a_real_action")
	}
	// Every real action keeps its default.
	if got := km.Stats.Help().Key; got != "ctrl+s" {
		t.Errorf("Stats unchanged default expected, got %q", got)
	}
}

func TestBuildKeyMap_InvalidKeyFallsBackToDefault(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionStats: "crtl+s", // typo
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 1 {
		t.Fatalf("issues len = %d, want 1. issues=%v", len(issues), issues)
	}
	if issues[0].Kind != IssueInvalidKey {
		t.Errorf("issue[0].Kind = %v, want IssueInvalidKey", issues[0].Kind)
	}
	if issues[0].Action != ActionStats {
		t.Errorf("issue[0].Action = %q, want %q", issues[0].Action, ActionStats)
	}
	if issues[0].Reason == "" {
		t.Errorf("issue[0].Reason should explain why the key was rejected")
	}
	// The user-facing message should include the reason.
	if msg := issues[0].Message(); !strings.Contains(msg, "crtl+s") || !strings.Contains(msg, issues[0].Reason) {
		t.Errorf("Message() = %q, want it to mention the bad key and the reason", msg)
	}
	if got := km.Stats.Help().Key; got != "ctrl+s" {
		t.Errorf("Stats falls back to default, got %q", got)
	}
}

func TestBuildKeyMap_DuplicateKey(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionStats:         "ctrl+d",
		ActionSchemaExplore: "ctrl+d",
	}
	_, issues := buildKeyMap(cfg)
	var dup *KeybindingIssue
	for i := range issues {
		if issues[i].Kind == IssueDuplicateKey {
			dup = &issues[i]
			break
		}
	}
	if dup == nil {
		t.Fatalf("expected IssueDuplicateKey, got %v", issues)
	}
	if dup.Key != "ctrl+d" {
		t.Errorf("duplicate key = %q, want %q", dup.Key, "ctrl+d")
	}
	acts := strings.Join(dup.Actions, ",")
	if !strings.Contains(acts, ActionStats) || !strings.Contains(acts, ActionSchemaExplore) {
		t.Errorf("duplicate actions = %v, want both %q and %q present", dup.Actions, ActionStats, ActionSchemaExplore)
	}
}

func TestBuildKeyMap_ActionForTranslatesOverrideAndDefault(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionStats: "ctrl+d",
	}
	km, _ := buildKeyMap(cfg)
	// Override key routes to the default.
	if got := km.translate("ctrl+d"); got != "ctrl+s" {
		t.Errorf("translate(override) = %q, want default ctrl+s", got)
	}
	// Non-rebindable keys pass through unchanged.
	if got := km.translate("tab"); got != "tab" {
		t.Errorf("translate(tab) = %q, want unchanged", got)
	}
	// Unset action's key still routes to itself (primary default).
	if got := km.translate("/"); got != "/" {
		t.Errorf("translate(/) = %q, want %q", got, "/")
	}
}

func TestRebindableActions_StableAndComplete(t *testing.T) {
	actions := rebindableActions()
	if len(actions) != len(defaultKeybindings) {
		t.Errorf("rebindableActions len = %d, want %d", len(actions), len(defaultKeybindings))
	}
	// Spot-check that well-known actions are present.
	want := []string{ActionStats, ActionSettings, ActionHelp, ActionQuit, ActionSearch}
	for _, w := range want {
		found := false
		for _, a := range actions {
			if a == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("action %q missing from rebindableActions", w)
		}
	}
}

func TestHelpEntries_ReflectsOverrides(t *testing.T) {
	cfg := &core.VaultConfig{}
	cfg.TUI.Keybindings = map[string]string{
		ActionStats: "ctrl+d",
	}
	km, _ := buildKeyMap(cfg)
	entries := helpEntries(km, false)
	var found *helpEntry
	for i, e := range entries {
		if e.Desc == "stats" {
			found = &entries[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("help entry for stats missing")
	}
	if found.Key != "ctrl+d" {
		t.Errorf("stats help key = %q, want %q after override", found.Key, "ctrl+d")
	}
}

func TestKeyMap_TranslateRoutesAliasToPrimary(t *testing.T) {
	// With defaults, pressing any bound key (primary OR alias) should
	// translate to the primary default — the switch in updateNormal uses
	// only primary keys as case labels.
	km := defaultKeyMap()
	for _, a := range rebindableActions() {
		primary := defaultKeybindings[a].Keys[0]
		for _, k := range defaultKeybindings[a].Keys {
			if got := km.translate(k); got != primary {
				t.Errorf("translate(%q) for action %q = %q, want primary %q", k, a, got, primary)
			}
		}
	}
}

func TestKeyMap_TranslateNonActionPassthrough(t *testing.T) {
	km := defaultKeyMap()
	// Random unbound keys pass through unchanged.
	for _, s := range []string{"x", "ctrl+x", "f5", "unknownkey"} {
		if got := km.translate(s); got != s {
			t.Errorf("translate(%q) = %q, want unchanged", s, got)
		}
	}
}

// Regression: tui-keybindings spec R1 S1 negative — rebind must act as a
// replace, not an add. After moving stats to ctrl+d, the old default ctrl+s
// must not dispatch to the stats case (issue #396).
func TestKeyMap_TranslateRebindUnbindsOldDefault(t *testing.T) {
	cfg := &core.VaultConfig{
		TUI: core.TUIConfig{
			Keybindings: map[string]string{
				ActionStats: "ctrl+d",
			},
		},
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 0 {
		t.Fatalf("unexpected validation issues: %+v", issues)
	}
	if got := km.translate("ctrl+d"); got != "ctrl+s" {
		t.Errorf("translate(ctrl+d) = %q, want ctrl+s (new key routes to primary)", got)
	}
	if got := km.translate("ctrl+s"); got == "ctrl+s" {
		t.Errorf("translate(ctrl+s) = %q — old default must NOT dispatch to stats case after rebind", got)
	}
	if got := km.translate("ctrl+s"); got != unboundSentinel {
		t.Errorf("translate(ctrl+s) = %q, want unboundSentinel — unbound default should be sentinel", got)
	}
}

// Regression: alias keys (e.g. `h` for help, which has primary `?`) should
// also unbind cleanly when the action is rebound. Every compile-time default
// key — primary or alias — must be neutralised if it is no longer bound.
func TestKeyMap_TranslateRebindUnbindsAliasDefaults(t *testing.T) {
	cfg := &core.VaultConfig{
		TUI: core.TUIConfig{
			Keybindings: map[string]string{
				ActionHelp: "f1",
			},
		},
	}
	km, issues := buildKeyMap(cfg)
	if len(issues) != 0 {
		t.Fatalf("unexpected validation issues: %+v", issues)
	}
	// Both `?` (primary) and `h` (alias) must no longer route to help.
	for _, oldKey := range []string{"?", "h"} {
		if got := km.translate(oldKey); got != unboundSentinel {
			t.Errorf("translate(%q) = %q, want unboundSentinel", oldKey, got)
		}
	}
	if got := km.translate("f1"); got != "?" {
		t.Errorf("translate(f1) = %q, want ? (new key routes to primary)", got)
	}
}

func TestConfigEditor_ListsKeybindingsInTUICategory(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init: %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open: %v", err)
	}
	defer v.Close()

	ce := newConfigEditor(v)
	// Locate the TUI category and check that every rebindable action's
	// "tui.keybindings.<action>" key is listed.
	var tuiCat *configCategory
	for i := range ce.categories {
		if ce.categories[i].Name == "TUI" {
			tuiCat = &ce.categories[i]
			break
		}
	}
	if tuiCat == nil {
		t.Fatalf("TUI category not found among %d categories", len(ce.categories))
	}
	keysInCat := map[string]bool{}
	for _, k := range tuiCat.Keys {
		keysInCat[k.Key] = true
	}
	for _, action := range rebindableActions() {
		want := "tui.keybindings." + action
		if !keysInCat[want] {
			t.Errorf("config editor TUI category missing key %q", want)
		}
	}
}

func TestConfigSetKeybinding_E2E(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init: %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open: %v", err)
	}
	defer v.Close()

	// tmd config set tui.keybindings.stats ctrl+d
	if err := v.SetConfigValue("tui.keybindings.stats", "ctrl+d"); err != nil {
		t.Fatalf("SetConfigValue: %v", err)
	}
	// tmd config get tui.keybindings.stats
	got, err := v.GetConfigValue("tui.keybindings.stats")
	if err != nil {
		t.Fatalf("GetConfigValue: %v", err)
	}
	if got != "ctrl+d" {
		t.Errorf("GetConfigValue = %q, want %q", got, "ctrl+d")
	}

	// Re-open the vault and confirm the value persisted to disk.
	v2 := core.NewVault(dir)
	if err := v2.Open(); err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer v2.Close()
	got2, err := v2.GetConfigValue("tui.keybindings.stats")
	if err != nil {
		t.Fatalf("second GetConfigValue: %v", err)
	}
	if got2 != "ctrl+d" {
		t.Errorf("after reopen, got %q, want %q", got2, "ctrl+d")
	}

	// And buildKeyMap with the reopened config returns the overridden key.
	km, issues := buildKeyMap(v2.Config())
	if len(issues) != 0 {
		t.Errorf("unexpected issues: %v", issues)
	}
	if helpKey := km.Stats.Help().Key; helpKey != "ctrl+d" {
		t.Errorf("Stats help key after reopen = %q, want %q", helpKey, "ctrl+d")
	}
}

func TestConfigKeysInfo_IncludesKeybindingEntries(t *testing.T) {
	// Use an unopened vault — ConfigKeysInfo reads the in-memory registry and
	// getter functions on the vault config, both of which are set up by
	// vault.Open. For a registry-only view we create a vault value without
	// opening and use the registered getters via ConfigKeys.
	keys := core.ConfigKeys()
	prefix := "tui.keybindings."
	var found []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			found = append(found, k)
		}
	}
	if len(found) != len(defaultKeybindings) {
		t.Errorf("registered keybinding keys = %d, want %d. found=%v", len(found), len(defaultKeybindings), found)
	}
	// Every rebindable action must have a matching registry entry.
	for _, a := range rebindableActions() {
		want := prefix + a
		hit := false
		for _, k := range keys {
			if k == want {
				hit = true
				break
			}
		}
		if !hit {
			t.Errorf("registry missing entry for %q", want)
		}
	}
	// Defaults should match the compile-time defaults.
	for _, a := range rebindableActions() {
		got := core.ConfigKeyDefault(prefix + a)
		want := defaultKeybindings[a].HelpKey
		if got != want {
			t.Errorf("ConfigKeyDefault(%s) = %q, want %q", prefix+a, got, want)
		}
	}
}

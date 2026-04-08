package core

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Domain steps
	dc := newDomainContext()

	// Cleanup after each scenario
	ctx.After(func(hookCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if dc.vault != nil {
			dc.vault.Close()
		}
		if dc.rootDir != "" {
			os.RemoveAll(dc.rootDir)
		}
		return hookCtx, nil
	})

	initCommonSteps(ctx, dc)
	initVaultSteps(ctx, dc)
	initObjectSteps(ctx, dc)
	initRelationSteps(ctx, dc)
	initQuerySteps(ctx, dc)
	initValidateSteps(ctx, dc)
	initWikiLinkSteps(ctx, dc)
	initResolveSteps(ctx, dc)
	initNameSteps(ctx, dc)
	initSharedSteps(ctx, dc)
	initSystemSteps(ctx, dc)
	initTagSteps(ctx, dc)
	initNameTemplateSteps(ctx, dc)
	initUniqueSteps(ctx, dc)
	initTemplateSteps(ctx, dc)
	initTypeCrudSteps(ctx, dc)
	initDoctorSteps(ctx, dc)
	initStarterSteps(ctx, dc)
	initVaultConfigSteps(ctx, dc)

	initConfigMgmtSteps(ctx, dc)
	initPageSteps(ctx, dc)
	initTypeDirectorySteps(ctx, dc)
	vcCtx := initViewConfigSteps(ctx, dc)
	initViewCrudSteps(ctx, dc, vcCtx)
	initQuerySortSteps(ctx, dc)
	initStatsSteps(ctx, dc)
	initRelationSyncSteps(ctx, dc)
	initFormatSteps(ctx, dc)
	initAIConfigSteps(ctx, dc)
	initAIAvailabilitySteps(ctx, dc)
	initAIMultiProviderSteps(ctx, dc)
	initLockSteps(ctx, dc)
	initArchiveSteps(ctx, dc)
	initDateDisplayFormatSteps(ctx, dc)

	initGraphSteps(ctx, dc)
	initMigrationSteps(ctx, dc)
	initImportScanSteps(ctx, dc)
	initImportPlanSteps(ctx, dc)
	initImportExecuteSteps(ctx, dc)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

package gormx_test

import "testing"

type dbScenario struct {
	name string
	test func(t *testing.T)
}

func gormxExecutionScenarios() []dbScenario {
	return []dbScenario{
		{name: "BaseRepo_SelectOneByOpts", test: TestBaseRepo_SelectOneByOpts},
		{name: "BaseRepo_SelectListByOpts", test: TestBaseRepo_SelectListByOpts},
		{name: "BaseRepo_SelectOneByMap", test: TestBaseRepo_SelectOneByMap},
		{name: "BaseRepo_SelectListByMap", test: TestBaseRepo_SelectListByMap},
		{name: "BaseRepo_SelectPage", test: TestBaseRepo_SelectPage},
		{name: "BaseRepo_SelectCount", test: TestBaseRepo_SelectCount},
		{name: "BaseRepo_Exists", test: TestBaseRepo_Exists},
		{name: "BaseRepo_Insert", test: TestBaseRepo_Insert},
		{name: "BaseRepo_InsertBatch", test: TestBaseRepo_InsertBatch},
		{name: "BaseRepo_InsertInBatches", test: TestBaseRepo_InsertInBatches},
		{name: "BaseRepo_InsertOrUpdate", test: TestBaseRepo_InsertOrUpdate},
		{name: "BaseRepo_Update", test: TestBaseRepo_Update},
		{name: "BaseRepo_UpdateById", test: TestBaseRepo_UpdateById},
		{name: "BaseRepo_UpdateByOpts", test: TestBaseRepo_UpdateByOpts},
		{name: "BaseRepo_Upsert", test: TestBaseRepo_Upsert},
		{name: "BaseRepo_GetOrCreate", test: TestBaseRepo_GetOrCreate},
		{name: "BaseRepo_UpdateOrCreate", test: TestBaseRepo_UpdateOrCreate},
		{name: "BaseRepo_Delete", test: TestBaseRepo_Delete},
		{name: "BaseRepo_DeleteById", test: TestBaseRepo_DeleteById},
		{name: "BaseRepo_DeleteBatchIds", test: TestBaseRepo_DeleteBatchIds},
		{name: "BaseRepo_DeleteByOpts", test: TestBaseRepo_DeleteByOpts},
		{name: "BaseRepo_DeleteByMap", test: TestBaseRepo_DeleteByMap},
		{name: "GenericAssociationManager_Add", test: TestGenericAssociationManager_Add},
		{name: "GenericAssociationManager_Remove", test: TestGenericAssociationManager_Remove},
		{name: "GenericAssociationManager_Clear", test: TestGenericAssociationManager_Clear},
		{name: "GenericAssociationManager_Set", test: TestGenericAssociationManager_Set},
		{name: "GenericAssociationManager_Count", test: TestGenericAssociationManager_Count},
		{name: "GenericAssociationManager_All", test: TestGenericAssociationManager_All},
		{name: "Query_First", test: TestQuery_First},
		{name: "Query_Find", test: TestQuery_Find},
		{name: "Query_Scan", test: TestQuery_Scan},
		{name: "Query_RawRows", test: TestQuery_RawRows},
		{name: "Query_Pluck", test: TestQuery_Pluck},
		{name: "Query_Preload", test: TestQuery_Preload},
	}
}

func runExecutionScenarios(t *testing.T, dialect string) {
	t.Helper()

	withTestDatabase(t, dialect, func(t *testing.T) {
		for _, scenario := range gormxExecutionScenarios() {
			t.Run(scenario.name, scenario.test)
		}
	})
}

func TestMySQLScenarios(t *testing.T) {
	runExecutionScenarios(t, testDialectMySQL)
}

func TestPostgreSQLScenarios(t *testing.T) {
	runExecutionScenarios(t, testDialectPostgres)
}

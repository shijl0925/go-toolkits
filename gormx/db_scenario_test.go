package gormx_test

import "testing"

type dbScenario struct {
	name string
	test func(t *testing.T)
}

func queryExecutionScenarios() []dbScenario {
	return []dbScenario{
		{name: "Query_Eq", test: TestQuery_Eq},
		{name: "Query_Ne", test: TestQuery_Ne},
		{name: "Query_Gt", test: TestQuery_Gt},
		{name: "Query_Ge", test: TestQuery_Ge},
		{name: "Query_Lt", test: TestQuery_Lt},
		{name: "Query_Le", test: TestQuery_Le},
		{name: "Query_Like", test: TestQuery_Like},
		{name: "Query_Regexp", test: TestQuery_Regexp},
		{name: "Query_IsNull", test: TestQuery_IsNull},
		{name: "Query_IsNotNull", test: TestQuery_IsNotNull},
		{name: "Query_In", test: TestQuery_In},
		{name: "Query_NotIn", test: TestQuery_NotIn},
		{name: "Query_Between", test: TestQuery_Between},
		{name: "Query_NotBetween", test: TestQuery_NotBetween},
		{name: "Query_OrderDesc", test: TestQuery_OrderDesc},
		{name: "Query_OrderAsc", test: TestQuery_OrderAsc},
		{name: "Query_First", test: TestQuery_First},
		{name: "Query_Find", test: TestQuery_Find},
		{name: "Query_Distinct", test: TestQuery_Distinct},
		{name: "Query_Select", test: TestQuery_Select},
		{name: "Query_Scan", test: TestQuery_Scan},
		{name: "Query_RawRows", test: TestQuery_RawRows},
		{name: "Query_Pluck", test: TestQuery_Pluck},
		{name: "Query_Preload", test: TestQuery_Preload},
		{name: "Query_GroupBy", test: TestQuery_GroupBy},
		{name: "Query_InSql", test: TestQuery_InSql},
		{name: "Query_NotInSql", test: TestQuery_NotInSql},
		{name: "Query_GtSql", test: TestQuery_GtSql},
		{name: "Query_GeSql", test: TestQuery_GeSql},
		{name: "Query_LtSql", test: TestQuery_LtSql},
		{name: "Query_LeSql", test: TestQuery_LeSql},
		{name: "Query_Not", test: TestQuery_Not},
		{name: "Query_Or", test: TestQuery_Or},
		{name: "Query_SubQueryEq", test: TestQuery_SubQueryEq},
		{name: "Query_SubQueryIn", test: TestQuery_SubQueryIn},
		{name: "Query_Count", test: TestQuery_Count},
		{name: "Query_Sum", test: TestQuery_Sum},
		{name: "Query_Avg", test: TestQuery_Avg},
		{name: "Query_Max", test: TestQuery_Max},
		{name: "Query_Min", test: TestQuery_Min},
		{name: "Query_Join", test: TestQuery_Join},
	}
}

func gormxExecutionScenarios() []dbScenario {
	scenarios := []dbScenario{
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
		{name: "BaseRepo_Transaction_Rollback", test: TestBaseRepo_Insert_WithTransaction_Rollback},
		{name: "BaseRepo_Transaction_Commit", test: TestBaseRepo_Insert_WithTransaction_Commit},
		{name: "BaseRepo_MultipleOps_Transaction", test: TestBaseRepo_MultipleOps_SameTransaction},
		{name: "BaseRepo_Context_Propagation", test: TestBaseRepo_UsesPassedContext},
		{name: "GenericAssociationManager_Add", test: TestGenericAssociationManager_Add},
		{name: "GenericAssociationManager_Remove", test: TestGenericAssociationManager_Remove},
		{name: "GenericAssociationManager_Clear", test: TestGenericAssociationManager_Clear},
		{name: "GenericAssociationManager_Set", test: TestGenericAssociationManager_Set},
		{name: "GenericAssociationManager_Count", test: TestGenericAssociationManager_Count},
		{name: "GenericAssociationManager_All", test: TestGenericAssociationManager_All},
	}

	return append(scenarios, queryExecutionScenarios()...)
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

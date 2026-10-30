SELECT t.* FROM app.f(1, x => 2) @{test_hint=1} AS t TABLESAMPLE RESERVOIR (2 ROWS)

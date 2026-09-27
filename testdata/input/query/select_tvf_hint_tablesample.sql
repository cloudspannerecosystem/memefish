SELECT * FROM app.f(1, x => 2) @{test_hint=1} TABLESAMPLE RESERVOIR (2 ROWS)

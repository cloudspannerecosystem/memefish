-- https://cloud.google.com/spanner/docs/ml-tutorial-generative-ai?hl=en#register_a_generative_ai_model_in_a_schema
SELECT content
FROM ML.PREDICT(
        MODEL TextBison,
        (SELECT "Is 13 prime?" AS prompt),
        STRUCT(256 AS maxOutputTokens, 0.2 AS temperature, 40 as topK, 0.95 AS topP)
     ) @{remote_udf_max_rows_per_rpc=1}
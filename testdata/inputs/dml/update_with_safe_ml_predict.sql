-- https://cloud.google.com/spanner/docs/backfill-embeddings?hl=en#backfill
UPDATE products
SET
    products.desc_embed = (
        SELECT embeddings.values
        FROM SAFE.ML.PREDICT(
                MODEL gecko_model,
                (SELECT products.description AS content)
             ) @{remote_udf_max_rows_per_rpc=200}
    ),
    products.desc_embed_model_version = 3
WHERE products.desc_embed IS NULL
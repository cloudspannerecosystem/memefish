CREATE MODEL MyClassificationModel
INPUT(
  content STRING(MAX)
)
OUTPUT(
  embeddings
    STRUCT<
      statistics STRUCT<truncated BOOL, token_count FLOAT64>,
      values ARRAY<FLOAT64>>
)
REMOTE
OPTIONS (
  endpoint = '//aiplatform.googleapis.com/projects/PROJECT/locations/LOCATION/endpoints/ENDPOINT_ID'
)
CREATE MODEL MyClassificationModel
INPUT (
  length FLOAT64,
  material STRING(MAX),
  tag_array ARRAY<STRING(MAX)>
)
OUTPUT (
  scores ARRAY<FLOAT64>,
  classes ARRAY<STRING(MAX)>
)
REMOTE
OPTIONS (
  endpoint = '//aiplatform.googleapis.com/projects/PROJECT/locations/LOCATION/endpoints/ENDPOINT_ID'
)
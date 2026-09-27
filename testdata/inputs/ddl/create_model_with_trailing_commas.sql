CREATE MODEL MyModel INPUT(
  prompt STRING(MAX),
 ) OUTPUT(
  content STRING(MAX),
 ) REMOTE OPTIONS (
  endpoint = '//aiplatform.googleapis.com/projects/PROJECT/locations/LOCATION/endpoints/ENDPOINT_ID',
  default_batch_size = 1
)
CREATE OR REPLACE MODEL GeminiPro
INPUT (prompt STRING(MAX))
OUTPUT (content STRING(MAX))
REMOTE OPTIONS (
  endpoint = '//aiplatform.googleapis.com/projects/fake-project/locations/asia-northeast1/publishers/google/models/gemini-pro',
  default_batch_size = 1
)
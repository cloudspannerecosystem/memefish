ALTER MODEL IF EXISTS MyClassificationModel
SET OPTIONS (
    endpoints = [
        '//aiplatform.googleapis.com/projects/aaa/locations/tl/endpoints/aaa',
        '//aiplatform.googleapis.com/projects/aaa/locations/tl/endpoints/bbb'
    ],
    default_batch_size = 100
)
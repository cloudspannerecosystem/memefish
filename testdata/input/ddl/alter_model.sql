ALTER MODEL MyClassificationModel
SET OPTIONS (
    endpoints = [
        '//aiplatform.googleapis.com/projects/aaa/locations/tl/endpoints/aaa',
        '//aiplatform.googleapis.com/projects/aaa/locations/tl/endpoints/bbb'
    ],
    default_batch_size = 100
)
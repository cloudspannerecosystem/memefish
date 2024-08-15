CREATE VECTOR INDEX hello_vector_index ON hello(embedding)
WHERE embedding IS NOT NULL
OPTIONS(distance_type = 'COSINE')
CREATE VECTOR INDEX IF NOT EXISTS hello_vector_index ON hello(embedding)
OPTIONS(distance_type = 'COSINE')
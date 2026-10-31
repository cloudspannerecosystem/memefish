CREATE VECTOR INDEX Singer_vector_index ON Singers(embedding, tenant_id, category)
OPTIONS(distance_type = 'COSINE')

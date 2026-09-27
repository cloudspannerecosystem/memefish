create vector index Singer_vector_index on Singers(embedding)
storing (genre)
where embedding is not null
options(distance_type = 'COSINE')
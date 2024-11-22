SELECT product_id, product_name, content
FROM ML.PREDICT(
        MODEL TextBison,
        (SELECT
             product.id as product_id,
             product.name as product_name,
             CONCAT("Is this product safe for infants?", "\n",
                    "Product Name: ", product.name, "\n",
                    "Category Name: ", category.name, "\n",
                    "Product Description:", product.description) AS prompt
         FROM
             Products AS product JOIN Categories AS category
                                      ON product.category_id = category.id),
        STRUCT(100 AS maxOutputTokens)
     ) @{remote_udf_max_rows_per_rpc=1}
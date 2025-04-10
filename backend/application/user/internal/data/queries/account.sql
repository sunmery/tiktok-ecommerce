-- name: GetFavorites :many
SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       -- 图片信息
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       -- 属性信息
       pa.attributes
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
         JOIN users.favorites uf ON p.id = uf.product_id
WHERE uf.user_id = @user_id::UUID
ORDER BY uf.created_at DESC
LIMIT @page_size::int OFFSET @page::int;

-- name: SetFavorites :exec
INSERT INTO users.favorites(user_id, product_id)
VALUES (@user_id, @product_id);

-- name: DeleteFavorites :exec
DELETE
FROM users.favorites
WHERE user_id = @user_id
  AND product_id = @product_id;

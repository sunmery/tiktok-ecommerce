SET SEARCH_PATH to merchant,products;

-- name: GetMerchantProducts :many
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.merchant_id = p.merchant_id) AS images,
       pa.attributes,
       (SELECT jsonb_build_object(
                       'id', a.id,
                       'old_status', a.old_status,
                       'new_status', a.new_status,
                       'reason', a.reason,
                       'created_at', a.created_at
               )
        FROM products.product_audits a
        WHERE a.merchant_id = p.merchant_id
        ORDER BY a.created_at DESC
        LIMIT 1)                              AS latest_audit
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.merchant_id = @merchant_id
  AND p.deleted_at IS NULL;

-- name: UpdateProduct :exec
UPDATE products.products
SET name        = coalesce(sqlc.narg(name), name),
    description = coalesce(sqlc.narg(description), description),
    price       = coalesce(sqlc.narg(price), price),
    updated_at  = now()
WHERE id = sqlc.arg(id)
  AND merchant_id = sqlc.arg(merchant_id);

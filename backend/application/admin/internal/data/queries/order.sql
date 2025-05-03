-- name: GetAllOrders :many
SELECT os.id     AS sub_order_id,
       os.total_amount,
       os.currency,
       os.status AS payment_status,
       os.items,
       os.shipping_status,
       os.created_at,
       os.updated_at,
       oo.id     AS order_id,
       oo.user_id,
       json_agg(
               json_build_object(
                       'streetAddress', oo.street_address,
                       'city', oo.city,
                       'state', oo.state,
                       'country', oo.country,
                       'zipCode', oo.zip_code
               )
       )         AS consumer_address,
       oo.email
FROM orders.sub_orders os
         JOIN orders.orders oo
              ON os.order_id = oo.id
group by oo.user_id, os.updated_at, os.id, os.id, os.total_amount, os.currency, os.status, os.items, os.shipping_status,
         os.created_at, os.updated_at, oo.id, oo.user_id, oo.email
ORDER BY os.created_at DESC
LIMIT @page_size OFFSET @page;

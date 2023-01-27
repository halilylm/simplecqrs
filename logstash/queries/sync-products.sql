SELECT l.id,
       l.operation,
       l.product_id,
       p.id,
       p.name,
       p.price
FROM products_logs l
         LEFT JOIN products p
                   ON p.id = l.product_id
WHERE l.id > :sql_last_value ORDER BY l.id;
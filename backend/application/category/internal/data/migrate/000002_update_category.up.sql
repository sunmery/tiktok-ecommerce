-- 修改分类表的 level 约束
-- ALTER TABLE categories.categories
-- DROP CONSTRAINT level_check;

ALTER TABLE categories.categories
ADD CONSTRAINT level_check CHECK (level BETWEEN 1 AND 4);

-- 修改闭包表的 depth 约束
-- ALTER TABLE categories.category_closure
-- DROP CONSTRAINT depth_check;

ALTER TABLE categories.category_closure
ADD CONSTRAINT depth_check CHECK (depth BETWEEN 0 AND 3); -- 0 到 3 表示四层
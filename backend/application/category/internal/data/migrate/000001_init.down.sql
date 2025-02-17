DROP TRIGGER IF EXISTS check_category_level_trigger ON categories.categories;
DROP FUNCTION IF EXISTS categories.check_category_level;
DROP TABLE IF EXISTS categories.category_closure;
DROP TABLE IF EXISTS categories.categories;
DROP SCHEMA IF EXISTS categories;

-- 宠物食物库存列（接管被删除的 000008_pet_food_inventory：当时与 000008_parent_isolation 版本号冲突，从未被执行）+
-- 已花费星星数（星星兑换零食）。
-- 注意：MySQL 8 不支持 ADD COLUMN IF NOT EXISTS，幂等性由代码层保证：
-- 迁移器执行前逐列检查 information_schema，列已存在则跳过对应语句。
ALTER TABLE pets ADD COLUMN food_apple INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN food_fish INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN food_milk INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN food_star INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN food_cake INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN food_hotpot INT NOT NULL DEFAULT 0;
ALTER TABLE pets ADD COLUMN stars_spent INT NOT NULL DEFAULT 0;

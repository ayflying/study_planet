-- 宠物表：记录已花费的星星数（用于星星兑换零食）
ALTER TABLE pets ADD COLUMN stars_spent INT NOT NULL DEFAULT 0 AFTER food_hotpot;
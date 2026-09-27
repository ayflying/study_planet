-- 000010 宠物衰减拆分为独立时间戳。
-- 原实现饱食度（24h 扣 20%）与好感度（8h -1）共用 last_decay_at，
-- 因好感度每 8h 衰减会不断刷新该时间戳，导致饱食度的 24h 衰减对频繁打开宠物页的孩子永不触发。
-- 拆分为 last_hunger_at / last_affection_at，各自按真实流逝时间独立结算。
ALTER TABLE pets ADD COLUMN last_hunger_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE pets ADD COLUMN last_affection_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

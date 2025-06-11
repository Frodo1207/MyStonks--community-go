-- 插入新手任务数据
INSERT INTO `tasks` (
    `id`, `step`, `title`, `description`, `reward`, `category`,
    `special_action`, `created_by`, `updated_by`
) VALUES
      (1, 1, '首次登录', '用你的web3钱包登录', 200, 'newbie',
       'login_popup', 'system', 'system'),
      (2, 2, '绑定TG', '绑定TG钱包', 100, 'newbie',
       'tg_bind_popup', 'system', 'system');

-- 插入日常任务数据
INSERT INTO `tasks` (
    `id`, `title`, `description`, `reward`, `category`,
    `icon`, `special_action`, `created_by`, `updated_by`
) VALUES
      (101, '每日签到', '访问社区TG并签到', 50, 'daily',
       'check-circle', NULL, 'system', 'system'),
      (102, '交易一次Stonks', '进行一次Stonks交易', 100, 'daily',
       'message-square', 'stonks_trade', 'system', 'system');

-- 插入其他任务数据
INSERT INTO `tasks` (
    `id`, `title`, `description`, `reward`, `category`,
    `created_by`, `updated_by`
) VALUES
    (201, '优质内容创作', '发布一篇被管理员标记为优质的内容', 500, 'other',
     'system', 'system');

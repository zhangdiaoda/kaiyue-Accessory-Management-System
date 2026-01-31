-- 权限管理与审计日志系统 - 数据库迁移脚本
-- 执行时间: 2026-01-31

-- 1. 权限表
CREATE TABLE IF NOT EXISTS `sys_permission` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(50) NOT NULL COMMENT '权限编码',
    `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `category` VARCHAR(50) NOT NULL COMMENT '权限分类',
    `description` VARCHAR(200) DEFAULT NULL COMMENT '权限描述',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统权限表';

-- 2. 用户权限关联表
CREATE TABLE IF NOT EXISTS `sys_user_permission` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `permission_code` VARCHAR(50) NOT NULL COMMENT '权限编码',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_permission` (`user_id`, `permission_code`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_permission_code` (`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限关联表';

-- 3. 角色默认权限表
CREATE TABLE IF NOT EXISTS `sys_role_permission` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role` VARCHAR(50) NOT NULL COMMENT '角色',
    `permission_code` VARCHAR(50) NOT NULL COMMENT '权限编码',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role`, `permission_code`),
    KEY `idx_role` (`role`),
    KEY `idx_permission_code` (`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色默认权限表';

-- 4. 操作日志表
CREATE TABLE IF NOT EXISTS `sys_operation_log` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '操作用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `real_name` VARCHAR(50) NOT NULL COMMENT '真实姓名',
    `operation` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `module` VARCHAR(50) NOT NULL COMMENT '模块',
    `description` TEXT COMMENT '操作描述',
    `request_method` VARCHAR(10) DEFAULT NULL COMMENT '请求方法',
    `request_url` VARCHAR(500) DEFAULT NULL COMMENT '请求URL',
    `request_params` TEXT COMMENT '请求参数',
    `response_result` TEXT COMMENT '响应结果',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    `status` VARCHAR(20) DEFAULT NULL COMMENT '操作状态: SUCCESS, FAILED',
    `error_message` TEXT COMMENT '错误信息',
    `execution_time` INT DEFAULT NULL COMMENT '执行时长(ms)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_operation` (`operation`),
    KEY `idx_module` (`module`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 初始化权限数据
INSERT INTO `sys_permission` (`code`, `name`, `category`, `description`) VALUES
-- 配件管理
('part:view', '查看配件', '配件管理', '查看配件列表和详情'),
('part:add', '添加配件', '配件管理', '添加新配件入库'),
('part:edit', '编辑配件', '配件管理', '修改配件信息'),
('part:delete', '删除配件', '配件管理', '删除配件记录'),
('part:import', '导入配件', '配件管理', '批量导入配件数据'),
('part:export', '导出配件', '配件管理', '导出配件数据'),

-- 领用管理
('borrow:create', '创建领用', '领用管理', '创建领用记录'),
('borrow:view', '查看领用记录', '领用管理', '查看领用记录'),
('borrow:return', '处理归还', '领用管理', '处理配件归还'),
('borrow:dispose', '处置旧件', '领用管理', '处置旧件和废品'),

-- 员工管理
('employee:view', '查看员工', '员工管理', '查看员工列表和详情'),
('employee:add', '添加员工', '员工管理', '添加新员工'),
('employee:edit', '编辑员工', '员工管理', '修改员工信息'),
('employee:delete', '删除员工', '员工管理', '删除员工记录'),
('employee:import', '导入员工', '员工管理', '批量导入员工数据'),
('employee:export', '导出员工', '员工管理', '导出员工数据'),

-- 报表管理
('report:view', '查看报表', '报表管理', '查看各类报表'),
('report:export', '导出报表', '报表管理', '导出报表数据');

-- 初始化仓库管理员默认权限
INSERT INTO `sys_role_permission` (`role`, `permission_code`) VALUES
('WAREHOUSE_ADMIN', 'part:view'),
('WAREHOUSE_ADMIN', 'part:add'),
('WAREHOUSE_ADMIN', 'part:edit'),
('WAREHOUSE_ADMIN', 'borrow:create'),
('WAREHOUSE_ADMIN', 'borrow:view'),
('WAREHOUSE_ADMIN', 'borrow:return'),
('WAREHOUSE_ADMIN', 'employee:view'),
('WAREHOUSE_ADMIN', 'report:view');

-- 创建数据库表结构
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 1. 用户表
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(100) NOT NULL COMMENT '密码（BCrypt加密）',
  `real_name` VARCHAR(50) NOT NULL COMMENT '真实姓名',
  `role` VARCHAR(20) NOT NULL COMMENT '角色：SUPER_ADMIN/WAREHOUSE_ADMIN',
  `department` VARCHAR(50) COMMENT '部门',
  `phone` VARCHAR(20) COMMENT '手机号',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0禁用/1启用',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_username` (`username`),
  INDEX `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';

-- ----------------------------
-- 2. 员工表
-- ----------------------------
DROP TABLE IF EXISTS `employee`;
CREATE TABLE `employee` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `employee_no` VARCHAR(50) NOT NULL UNIQUE COMMENT '员工工号',
  `name` VARCHAR(50) NOT NULL COMMENT '姓名',
  `department` VARCHAR(50) COMMENT '部门',
  `position` VARCHAR(50) COMMENT '岗位',
  `phone` VARCHAR(20) COMMENT '手机号',
  `status` TINYINT DEFAULT 1 COMMENT '状态：0离职/1在职',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_employee_no` (`employee_no`),
  INDEX `idx_department` (`department`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工表';

-- ----------------------------
-- 3. 配件分类表
-- ----------------------------
DROP TABLE IF EXISTS `part_category`;
CREATE TABLE `part_category` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父分类ID（0为顶级）',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配件分类表';

-- ----------------------------
-- 4. 配件信息表
-- ----------------------------
DROP TABLE IF EXISTS `part`;
CREATE TABLE `part` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `part_no` VARCHAR(50) NOT NULL UNIQUE COMMENT '配件编号',
  `name` VARCHAR(100) NOT NULL COMMENT '配件名称',
  `category_id` BIGINT NOT NULL COMMENT '分类ID',
  `specification` VARCHAR(200) COMMENT '规格型号',
  `unit` VARCHAR(20) DEFAULT '件' COMMENT '单位',
  `stock_quantity` INT DEFAULT 0 COMMENT '当前库存数量',
  `warning_threshold` INT DEFAULT 10 COMMENT '预警阈值',
  `price` DECIMAL(10, 2) COMMENT '单价',
  `remark` TEXT COMMENT '备注',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_part_no` (`part_no`),
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_stock_warning` (`stock_quantity`, `warning_threshold`),
  FOREIGN KEY (`category_id`) REFERENCES `part_category`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配件信息表';

-- ----------------------------
-- 5. 领用记录表
-- ----------------------------
DROP TABLE IF EXISTS `borrow_record`;
CREATE TABLE `borrow_record` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `record_no` VARCHAR(50) NOT NULL UNIQUE COMMENT '记录编号',
  `employee_id` BIGINT NOT NULL COMMENT '员工ID',
  `part_id` BIGINT NOT NULL COMMENT '配件ID',
  `borrow_quantity` INT NOT NULL COMMENT '领用数量',
  `return_quantity` INT DEFAULT 0 COMMENT '已归还数量',
  `damaged_quantity` INT DEFAULT 0 COMMENT '损毁数量',
  `status` VARCHAR(20) NOT NULL COMMENT '状态：BORROWED/RETURNED/PARTIAL_RETURNED',
  `borrow_time` DATETIME NOT NULL COMMENT '领用时间',
  `borrow_admin_id` BIGINT NOT NULL COMMENT '登记管理员ID',
  `return_time` DATETIME COMMENT '归还时间',
  `return_admin_id` BIGINT COMMENT '归还登记管理员ID',
  `remark` TEXT COMMENT '备注',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_employee_id` (`employee_id`),
  INDEX `idx_part_id` (`part_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_borrow_time` (`borrow_time`),
  FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`),
  FOREIGN KEY (`part_id`) REFERENCES `part`(`id`),
  FOREIGN KEY (`borrow_admin_id`) REFERENCES `sys_user`(`id`),
  FOREIGN KEY (`return_admin_id`) REFERENCES `sys_user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='领用记录表';

-- ----------------------------
-- 6. 系统配置表
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `config_key` VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
  `config_value` TEXT COMMENT '配置值',
  `description` VARCHAR(200) COMMENT '说明',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 插入默认管理员账号（密码：admin123）
-- BCrypt哈希生成：admin123
INSERT INTO `sys_user` (`username`, `password`, `real_name`, `role`, `department`, `phone`, `status`) VALUES
('admin', '$2a$10$rQ9P6zHxfKJWx0p3VZ9vL.JQvH3fHZJX6JN2L0ZqGOPJ1YJnZxZmK', '系统管理员', 'SUPER_ADMIN', '技术部', '13800138000', 1),
('warehouse', '$2a$10$rQ9P6zHxfKJWx0p3VZ9vL.JQvH3fHZJX6JN2L0ZqGOPJ1YJnZxZmK', '仓库管理员', 'WAREHOUSE_ADMIN', '仓储部', '13800138001', 1);

-- 插入默认配置
INSERT INTO `sys_config` (`config_key`, `config_value`, `description`) VALUES
('dingtalk.webhook.url', '', '钉钉群机器人Webhook地址'),
('dingtalk.secret', '', '钉钉群机器人加签密钥'),
('report.auto_push', 'true', '是否自动推送报表'),
('report.push_day', '1', '每月推送日期（1-28）'),
('stock.warning.enabled', 'true', '是否启用库存预警');

-- 插入默认配件分类
INSERT INTO `part_category` (`name`, `parent_id`, `sort_order`) VALUES
('刀具类', 0, 1),
('刀片', 1, 1),
('钻头', 1, 2),
('铣刀', 1, 3),
('测量工具', 0, 2),
('卡尺', 4, 1),
('千分尺', 4, 2),
('辅助工具', 0, 3);

-- 插入示例员工
INSERT INTO `employee` (`employee_no`, `name`, `department`, `position`, `phone`, `status`) VALUES
('EMP001', '张三', '生产部', '机械工程师', '13900139001', 1),
('EMP002', '李四', '生产部', '技术员', '13900139002', 1),
('EMP003', '王五', '质检部', '质检员', '13900139003', 1);

-- 插入示例配件
INSERT INTO `part` (`part_no`, `name`, `category_id`, `specification`, `unit`, `stock_quantity`, `warning_threshold`, `price`) VALUES
('PART001', '硬质合金刀片', 2, 'CNMG120408', '片', 100, 20, 15.50),
('PART002', '高速钢钻头', 3, 'Φ10mm', '支', 50, 10, 8.00),
('PART003', '立铣刀', 4, 'Φ12mm', '支', 30, 5, 25.00),
('PART004', '游标卡尺', 5, '0-150mm', '把', 15, 3, 45.00),
('PART005', '千分尺', 6, '0-25mm', '个', 8, 2, 120.00);

SET FOREIGN_KEY_CHECKS = 1;

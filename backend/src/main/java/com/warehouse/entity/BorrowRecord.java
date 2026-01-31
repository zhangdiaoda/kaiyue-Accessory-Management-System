package com.warehouse.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;

import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * 领用记录实体
 */
@Data
@TableName("borrow_record")
public class BorrowRecord implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    /**
     * 记录编号
     */
    private String recordNo;

    /**
     * 员工ID
     */
    private Long employeeId;

    /**
     * 配件ID
     */
    private Long partId;

    /**
     * 领用数量
     */
    private Integer borrowQuantity;

    /**
     * 已归还数量
     */
    private Integer returnQuantity;

    /**
     * 损毁数量
     */
    private Integer damagedQuantity;

    /**
     * 状态：BORROWED/RETURNED/PARTIAL_RETURNED
     */
    private String status;

    /**
     * 领用时间
     */
    private LocalDateTime borrowTime;

    /**
     * 登记管理员ID
     */
    private Long borrowAdminId;

    /**
     * 归还时间
     */
    private LocalDateTime returnTime;

    /**
     * 归还登记管理员ID
     */
    private Long returnAdminId;

    /**
     * 备注
     */
    private String remark;

    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createdAt;

    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updatedAt;
}

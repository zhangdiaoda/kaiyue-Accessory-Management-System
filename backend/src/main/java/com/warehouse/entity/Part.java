package com.warehouse.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 配件实体
 */
@Data
@TableName("part")
public class Part implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    /**
     * 配件编号
     */
    private String partNo;

    /**
     * 配件名称
     */
    private String name;

    /**
     * 分类ID
     */
    private Long categoryId;

    /**
     * 规格型号
     */
    private String specification;

    /**
     * 单位
     */
    private String unit;

    /**
     * 当前库存数量
     */
    private Integer stockQuantity;

    /**
     * 预警阈值
     */
    private Integer warningThreshold;

    /**
     * 单价
     */
    private BigDecimal price;

    /**
     * 备注
     */
    private String remark;

    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createdAt;

    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updatedAt;
}

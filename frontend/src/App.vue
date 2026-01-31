<template>
  <router-view />
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

/* ========== 移动端全局样式 ========== */
@media (max-width: 768px) {
  /* 字体缩放 */
  html {
    font-size: 14px;
  }
  
  /* 按钮最小触摸尺寸 44x44px (Apple HIG) */
  .el-button {
    min-height: 44px;
    padding: 0 20px;
  }
  
  .el-button--large {
    min-height: 48px;
  }
  
  .el-button--small {
    min-height: 40px;
  }
  
  /* 输入框高度 */
  .el-input__wrapper {
    min-height: 44px;
  }
  
  .el-textarea__inner {
    font-size: 14px;
  }
  
  /* 弹窗移动端适配 */
  .el-dialog {
    width: 95% !important;
    max-width: 500px;
    margin: 20px auto !important;
  }
  
  .el-dialog__header {
    padding: 16px 20px;
  }
  
  .el-dialog__body {
    padding: 20px;
  }
  
  .el-dialog__footer {
    padding: 12px 20px;
  }
  
  /* 表单间距 */
  .el-form-item {
    margin-bottom: 20px;
  }
  
  /* Select 下拉框 */
  .el-select {
    width: 100%;
  }
  
  /* 卡片间距 */
  .el-card {
    margin-bottom: 16px;
  }
  
  /* 分页器 */
  .el-pagination {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .el-pagination button,
  .el-pagination .el-pager li {
    min-width: 36px;
    height: 36px;
  }
}

/* 平板优化 */
@media (min-width: 768px) and (max-width: 1024px) {
  .el-dialog {
    width: 80% !important;
    max-width: 700px;
  }
}
</style>

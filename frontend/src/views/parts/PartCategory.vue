<template>
  <div class="part-category">
    <!-- Header Actions -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">配件分类库</h2>
        <span class="page-subtitle">管理所有资产的逻辑归属</span>
      </div>
      <el-button type="primary" size="large" class="apple-primary-btn" @click="handleAdd">
        <el-icon class="el-icon--left"><Plus /></el-icon>
        新建顶级分类
      </el-button>
    </div>

    <!-- Category Grid -->
    <el-row :gutter="24" class="category-grid">
      <!-- 顶级分类卡片 -->
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="(root, index) in tableData" :key="root.id">
        <div class="category-card animate-scale-in" :style="{ '--card-accent': getAccentColor(index) }">
          <div class="card-header-bar">
            <div class="cat-info">
              <span class="cat-name">{{ root.name }}</span>
              <el-tag size="small" effect="plain" round class="count-tag">{{ root.children?.length || 0 }} 子类</el-tag>
            </div>
            <div class="cat-actions">
              <el-dropdown trigger="click">
                <el-button link class="action-trigger">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu class="apple-dropdown">
                    <el-dropdown-item @click="handleAddChild(root)"><el-icon><Plus /></el-icon>添加子类</el-dropdown-item>
                    <el-dropdown-item @click="handleEdit(root)"><el-icon><Edit /></el-icon>编辑分类</el-dropdown-item>
                    <el-dropdown-item divided @click="handleDelete(root)" class="text-danger"><el-icon><Delete /></el-icon>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>

          <!-- 子分类平铺区 -->
          <div class="card-body">
            <div class="sub-tags-container" v-if="root.children && root.children.length > 0">
              <div 
                v-for="child in root.children" 
                :key="child.id" 
                class="sub-cat-chip"
                @click="handleEdit(child)"
              >
                <span class="sub-name">{{ child.name }}</span>
                <el-icon class="sub-arrow"><ArrowRight /></el-icon>
              </div>
            </div>
            <div class="empty-state" v-else>
              <span>暂无子分类</span>
              <el-button link type="primary" size="small" @click="handleAddChild(root)">立即添加</el-button>
            </div>
          </div>

          <!-- 底部装饰条 -->
          <div class="card-footer-deco"></div>
        </div>
      </el-col>
    </el-row>

    <!-- 新增/编辑对话框 -->
    <el-dialog 
      v-model="showDialog" 
      :title="dialogTitle" 
      width="420px"
      class="apple-dialog"
      append-to-body
    >
      <el-form :model="formData" label-position="top" class="apple-form">
        <el-form-item label="所属父级">
          <el-select 
            v-model="formData.parent_id" 
            placeholder="设为顶级分类" 
            clearable
            class="apple-select"
          >
            <el-option label="✦ 设为顶级分类" :value="0" />
            <el-option
              v-for="cat in flatCategories.filter(c => c.id !== formData.id)"
              :key="cat.id"
              :label="cat.label"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="分类名称" required>
          <el-input v-model="formData.name" placeholder="如: 精密刀具" class="apple-input" />
        </el-form-item>
        <el-form-item label="排序权重">
          <el-input-number v-model="formData.sort_order" :min="0" controls-position="right" style="width: 100%" />
          <div class="field-tips">数字越小排序越靠前</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showDialog = false" round>取消</el-button>
          <el-button type="primary" @click="handleSubmit" round class="apple-primary-btn">保存变更</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { MoreFilled, Plus, Edit, Delete, ArrowRight } from '@element-plus/icons-vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api/category'

const tableData = ref([])
const showDialog = ref(false)
const dialogTitle = ref('添加分类')

const formData = ref({
  id: null,
  name: '',
  parent_id: 0,
  sort_order: 0
})

// Apple 风格柔和色板
const accentColors = [
  '#0071e3', '#ff375f', '#ff9f0a', '#34c759', '#5e5ce6', '#af52de', '#5ac8fa'
]

const getAccentColor = (index) => {
  return accentColors[index % accentColors.length]
}

// 扁平化分类列表（用于下拉选择）
const flatCategories = computed(() => {
  const result = []
  const flatten = (items, level = 0) => {
    items.forEach(item => {
      result.push({
        id: item.id,
        label: '  '.repeat(level) + item.name,
        level
      })
      // 我们只允许两级分类，所以在 Dialog 中只显示顶级分类供选择
      // 如果后端支持无限层级，这里可以放开
    })
  }
  // 只取顶级分类作为父级选项
  flatten(tableData.value.filter(i => i.parent_id === 0))
  return result
})

const loadData = async () => {
  try {
    const res = await getCategoryList()
    tableData.value = buildTree(res.data || [])
  } catch (error) {
    ElMessage.error('加载分类数据失败')
  }
}

// 构建树形结构
const buildTree = (list) => {
  const map = {}
  const roots = []
  
  // 深拷贝以避免引用问题
  const rawList = JSON.parse(JSON.stringify(list))
  
  rawList.forEach(item => {
    map[item.id] = { ...item, children: [] }
  })
  
  rawList.forEach(item => {
    if (item.parent_id === 0) {
      roots.push(map[item.id])
    } else if (map[item.parent_id]) {
      map[item.parent_id].children.push(map[item.id])
    }
  })
  
  // 排序
  roots.sort((a, b) => a.sort_order - b.sort_order)
  roots.forEach(r => r.children.sort((a, b) => a.sort_order - b.sort_order))
  
  return roots
}

const handleAdd = () => {
  dialogTitle.value = '新建顶级分类'
  formData.value = {
    id: null,
    name: '',
    parent_id: 0,
    sort_order: 0
  }
  showDialog.value = true
}

const handleAddChild = (row) => {
  dialogTitle.value = `添加 "${row.name}" 的子类`
  formData.value = {
    id: null,
    name: '',
    parent_id: row.id,
    sort_order: 0
  }
  showDialog.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑分类'
  formData.value = { ...row }
  showDialog.value = true
}

const handleDelete = async (row) => {
  try {
    const hasChildren = row.children && row.children.length > 0
    if (hasChildren) {
      ElMessage.warning('请先删除该分类下的所有子分类')
      return
    }

    await ElMessageBox.confirm(
      `确定要删除分类 "${row.name}" 吗？此操作不可恢复。`, 
      '删除确认', 
      { 
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    await deleteCategory(row.id)
    ElMessage.success('已删除')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formData.value.name) {
    ElMessage.warning('请输入分类名称')
    return
  }

  try {
    if (formData.value.id) {
      await updateCategory(formData.value.id, formData.value)
      ElMessage.success('分类已更新')
    } else {
      await createCategory(formData.value)
      ElMessage.success('分类已创建')
    }
    showDialog.value = false
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.part-category {
  padding-bottom: 40px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 30px;
}

.page-title {
  font-size: 28px;
  font-weight: 800;
  color: #1d1d1f;
  margin: 0 0 4px 0;
  letter-spacing: -0.5px;
}

.page-subtitle {
  font-size: 14px;
  color: #86868b;
}

.category-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  border: 1px solid rgba(0,0,0,0.03);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 220px;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0,0,0,0.08);
}

.category-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: var(--card-accent);
  opacity: 0.8;
}

.card-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.cat-name {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #1d1d1f;
  margin-bottom: 6px;
}

.count-tag {
  background: rgba(0,0,0,0.03);
  border: none;
  font-weight: 600;
  color: #86868b;
}

.action-trigger {
  color: #86868b;
  font-size: 18px;
  padding: 4px;
}

.action-trigger:hover {
  color: var(--card-accent);
  background: rgba(0,0,0,0.03);
  border-radius: 8px;
}

.card-body {
  flex: 1;
}

.sub-tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.sub-cat-chip {
  background: #f5f5f7;
  padding: 8px 14px;
  border-radius: 12px;
  font-size: 13px;
  color: #1d1d1f;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.sub-cat-chip:hover {
  background: var(--card-accent);
  color: white;
}

.sub-arrow {
  opacity: 0;
  font-size: 12px;
  transform: translateX(-4px);
  transition: all 0.2s;
}

.sub-cat-chip:hover .sub-arrow {
  opacity: 1;
  transform: translateX(0);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #98989d;
  font-size: 13px;
  gap: 8px;
  opacity: 0.6;
}

.apple-dialog {
  border-radius: 20px;
}

.field-tips {
  font-size: 12px;
  color: #86868b;
  margin-top: 6px;
}

.text-danger {
  color: #ff375f;
}

/* Animations */
.animate-scale-in {
  animation: scaleIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
  transform: scale(0.95);
}

@keyframes scaleIn {
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* Delay for staggering effect */
.category-grid .el-col:nth-child(1) .category-card { animation-delay: 0.05s; }
.category-grid .el-col:nth-child(2) .category-card { animation-delay: 0.1s; }
.category-grid .el-col:nth-child(3) .category-card { animation-delay: 0.15s; }
.category-grid .el-col:nth-child(4) .category-card { animation-delay: 0.2s; }
</style>


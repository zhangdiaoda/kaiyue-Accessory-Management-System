<template>
  <el-dialog
    v-model="visible"
    :title="`设置权限 - ${user?.real_name || ''}`"
    width="700px"
    @close="handleClose"
  >
    <div v-loading="loading">
      <!-- 角色默认权限提示 -->
      <el-alert
        title="角色默认权限已灰色显示，无需勾选即自动拥有"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <!-- 权限列表（按分类分组） -->
      <div v-for="(perms, category) in permissionsByCategory" :key="category" class="permission-category">
        <div class="category-title">{{ category }}</div>
        <el-checkbox-group v-model="selectedPermissions" class="permission-list">
          <el-checkbox
            v-for="perm in perms"
            :key="perm.code"
            :label="perm.code"
            :disabled="isRolePermission(perm.code)"
          >
            {{ perm.name }}
            <el-tag v-if="isRolePermission(perm.code)" size="small" type="info" style="margin-left: 8px">
              默认
            </el-tag>
          </el-checkbox>
        </el-checkbox-group>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getAllPermissions, getUserPermissions, setUserPermissions, getRolePermissions } from '@/api/permission'

const props = defineProps({
  modelValue: Boolean,
  user: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const saving = ref(false)
const permissionsByCategory = ref({})
const selectedPermissions = ref([])
const roleDefaultPermissions = ref([])

// 判断是否为角色默认权限
const isRolePermission = (code) => {
  return roleDefaultPermissions.value.includes(code)
}

// 加载权限数据
const loadPermissions = async () => {
  if (!props.user) return

  loading.value = true
  try {
    // 1. 获取所有权限
    const allRes = await getAllPermissions()
    if (allRes.code === 200) {
      permissionsByCategory.value = allRes.data
    }

    // 2. 获取角色默认权限
    const roleRes = await getRolePermissions(props.user.role)
    if (roleRes.code === 200) {
      roleDefaultPermissions.value = roleRes.data || []
    }

    // 3. 获取用户当前权限（仅个人扩展权限，不包含角色默认权限）
    const userRes = await getUserPermissions(props.user.id)
    if (userRes.code === 200) {
      const allUserPerms = userRes.data || []
      // 过滤出非角色默认权限（个人扩展权限）
      selectedPermissions.value = allUserPerms.filter(
        code => !roleDefaultPermissions.value.includes(code)
      )
    }
  } catch (error) {
    ElMessage.error('加载权限失败')
  } finally {
    loading.value = false
  }
}

// 保存权限
const handleSave = async () => {
  saving.value = true
  try {
    // 只保存用户的个人扩展权限（不包含角色默认权限）
    await setUserPermissions(props.user.id, selectedPermissions.value)
    ElMessage.success('权限设置成功')
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error('权限设置失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  visible.value = false
  selectedPermissions.value = []
  roleDefaultPermissions.value = []
}

watch(() => props.modelValue, (val) => {
  if (val) {
    loadPermissions()
  }
})
</script>

<style scoped>
.permission-category {
  margin-bottom: 24px;
}

.category-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
}

.permission-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.permission-list :deep(.el-checkbox) {
  margin-right: 0;
}

.permission-list :deep(.el-checkbox.is-disabled) {
  opacity: 0.5;
}
</style>

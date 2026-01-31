<template>
  <div class="user-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加用户
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" border stripe>
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="role" label="角色" width="130">
          <template #default="{ row }">
            <el-tag :type="row.role === 'SUPER_ADMIN' ? 'danger' : 'primary'">
              {{ row.role === 'SUPER_ADMIN' ? '超级管理员' : '仓库管理员' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="department" label="部门" width="150" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.role === 'WAREHOUSE_ADMIN'" 
              link 
              type="success" 
              @click="handlePermission(row)"
            >
              权限设置
            </el-button>
            <el-button link type="warning" @click="handleResetPassword(row)">重置密码</el-button>
            <el-button v-if="row.username !== 'admin'" link type="danger" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="formData.username" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item v-if="!formData.id" label="密码" required>
          <el-input v-model="formData.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="真实姓名" required>
          <el-input v-model="formData.real_name" />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-radio-group v-model="formData.role">
            <el-radio label="SUPER_ADMIN">超级管理员</el-radio>
            <el-radio label="WAREHOUSE_ADMIN">仓库管理员</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="部门">
          <el-input v-model="formData.department" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限设置对话框 -->
    <PermissionDialog v-model="showPermissionDialog" :user="currentUser" @success="loadData" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getAllUsers, createUser, updateUser, deleteUser, resetUserPassword } from '@/api/auth'
import PermissionDialog from '@/components/PermissionDialog.vue'

const tableData = ref([])
const showDialog = ref(false)
const showPermissionDialog = ref(false)
const currentUser = ref(null)
const dialogTitle = ref('添加用户')

const formData = reactive({
  id: null,
  username: '',
  password: '',
  real_name: '',
  role: 'WAREHOUSE_ADMIN',
  department: '',
  phone: '',
  status: 1
})

const loadData = async () => {
  try {
    const res = await getAllUsers()
    if (res.code === 200) {
      tableData.value = res.data
    }
  } catch (error) {
    ElMessage.error('无法加载用户列表')
  }
}

const handleAdd = () => {
  dialogTitle.value = '添加用户'
  Object.assign(formData, {
    id: null,
    username: '',
    password: '',
    real_name: '',
    role: 'WAREHOUSE_ADMIN',
    department: '',
    phone: '',
    status: 1
  })
  showDialog.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑用户'
  Object.assign(formData, row)
  formData.password = '' // 编辑时不显示密码
  showDialog.value = true
}

const handleResetPassword = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入新密码', '重置密码', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'password',
      inputPattern: /.{6,}/,
      inputErrorMessage: '密码至少6位'
    })
    
    const res = await resetUserPassword(row.id, value)
    if (res.code === 200) {
      ElMessage.success('密码重置成功')
    }
  } catch (error) {
    // 用户取消或失败
  }
}

const handlePermission = (row) => {
  currentUser.value = row
  showPermissionDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗？', '提示', { type: 'warning' })
    const res = await deleteUser(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadData()
    } else {
      ElMessage.error(res.message)
    }
  } catch (error) {
    // 用户取消
  }
}

const handleSubmit = async () => {
  if (!formData.username || !formData.real_name) {
    ElMessage.warning('请填写必填项')
    return
  }
  if (!formData.id && !formData.password) {
    ElMessage.warning('请输入密码')
    return
  }

  try {
    let res
    if (formData.id) {
      res = await updateUser(formData.id, formData)
    } else {
      res = await createUser(formData)
    }

    if (res.code === 200) {
      ElMessage.success(formData.id ? '更新成功' : '创建成功')
      showDialog.value = false
      loadData()
    } else {
      ElMessage.error(res.message)
    }
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.user-manage {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

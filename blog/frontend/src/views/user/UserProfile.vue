<template>
  <div class="user-profile">
    <div class="page-header">
      <h2 class="page-title">个人资料</h2>
      <el-button type="primary" @click="isEditing = true" v-if="!isEditing">
        <el-icon><Edit /></el-icon>
        编辑资料
      </el-button>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton animated>
        <template #template>
          <div class="profile-skeleton">
            <div class="skeleton-item avatar-skeleton"></div>
            <div style="padding: 14px;">
              <div class="skeleton-item skeleton-title"></div>
              <div class="skeleton-item" style="margin-top: 16px; height: 16px;"></div>
              <div class="skeleton-item" style="margin-top: 16px; height: 16px;"></div>
              <div class="skeleton-item" style="margin-top: 16px; height: 16px;"></div>
            </div>
          </div>
        </template>
      </el-skeleton>
    </div>
    
    <div v-else-if="!isEditing" class="profile-display">
      <div class="profile-header">
        <div class="profile-avatar">
          <img :src="user.avatar || defaultAvatar" :alt="user.username" />
        </div>
        <div class="profile-basic">
          <h3 class="profile-name">{{ user.username }}</h3>
          <div class="profile-meta">
            <div class="meta-item">
              <el-icon><Message /></el-icon>
              {{ user.email }}
            </div>
            <div class="meta-item">
              <el-icon><Calendar /></el-icon>
              注册于 {{ formatDate(user.createdAt) }}
            </div>
            <div class="meta-item">
              <el-icon><Document /></el-icon>
              {{ user.postCount || 0 }} 篇文章
            </div>
            <div class="meta-item">
              <el-icon><ChatLineRound /></el-icon>
              {{ user.commentCount || 0 }} 条评论
            </div>
          </div>
        </div>
      </div>
      
      <div class="profile-content">
        <div class="profile-section">
          <h4 class="section-title">个人简介</h4>
          <div class="section-content bio-content">
            {{ user.bio || '这个人很懒，什么都没写...' }}
          </div>
        </div>
        
        <div class="profile-section">
          <h4 class="section-title">个人链接</h4>
          <div class="section-content">
            <div class="link-item" v-if="user.website">
              <el-icon><Link /></el-icon>
              <a :href="user.website" target="_blank" rel="noopener noreferrer">{{ user.website }}</a>
            </div>
            <div class="link-item" v-if="user.github">
              <svg viewBox="0 0 24 24" width="16" height="16" class="custom-icon">
                <path fill="currentColor" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.022A9.606 9.606 0 0112 6.82c.85.004 1.705.114 2.504.336 1.909-1.291 2.747-1.022 2.747-1.022.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.841-2.337 4.687-4.565 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.137 20.164 22 16.417 22 12c0-5.523-4.477-10-10-10z"/>
              </svg>
              <a :href="`https://github.com/${user.github}`" target="_blank" rel="noopener noreferrer">{{ user.github }}</a>
            </div>
            <div class="link-item" v-if="user.twitter">
              <svg viewBox="0 0 24 24" width="16" height="16" class="custom-icon">
                <path fill="currentColor" d="M22.46 6c-.77.35-1.6.58-2.46.69.88-.53 1.56-1.37 1.88-2.38-.83.5-1.75.85-2.72 1.05C18.37 4.5 17.26 4 16 4c-2.35 0-4.27 1.92-4.27 4.29 0 .34.04.67.11.98C8.28 9.09 5.11 7.38 3 4.79c-.37.63-.58 1.37-.58 2.15 0 1.49.75 2.81 1.91 3.56-.71 0-1.37-.2-1.95-.5v.03c0 2.08 1.48 3.82 3.44 4.21a4.22 4.22 0 0 1-1.93.07 4.28 4.28 0 0 0 4 2.98 8.521 8.521 0 0 1-5.33 1.84c-.34 0-.68-.02-1.02-.06C3.44 20.29 5.7 21 8.12 21 16 21 20.33 14.46 20.33 8.79c0-.19 0-.37-.01-.56.84-.6 1.56-1.36 2.14-2.23z"/>
              </svg>
              <a :href="`https://twitter.com/${user.twitter}`" target="_blank" rel="noopener noreferrer">{{ user.twitter }}</a>
            </div>
            <div class="empty-links" v-if="!user.website && !user.github && !user.twitter">
              暂无个人链接
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="profile-edit">
      <el-form :model="profileForm" :rules="rules" ref="profileFormRef" label-width="80px">
        <el-form-item label="头像">
          <div class="avatar-uploader">
            <img :src="avatarUrl || defaultAvatar" class="avatar-preview" />
            <div class="upload-actions">
              <el-button type="primary" @click="triggerUpload">
                <el-icon><Upload /></el-icon>
                上传头像
              </el-button>
              <input 
                type="file" 
                ref="fileInput" 
                style="display: none" 
                accept="image/*"
                @change="handleFileChange" 
              />
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="用户名" prop="username">
          <el-input v-model="profileForm.username" placeholder="请输入用户名" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" placeholder="请输入邮箱" disabled />
        </el-form-item>
        
        <el-form-item label="个人简介" prop="bio">
          <el-input
            v-model="profileForm.bio"
            type="textarea"
            :rows="4"
            placeholder="介绍一下自己吧..."
          />
        </el-form-item>
        
        <el-form-item label="个人网站" prop="website">
          <el-input v-model="profileForm.website" placeholder="个人网站地址 (选填)" />
        </el-form-item>
        
        <el-form-item label="GitHub" prop="github">
          <el-input v-model="profileForm.github" placeholder="GitHub 用户名 (选填)" />
        </el-form-item>
        
        <el-form-item label="Twitter" prop="twitter">
          <el-input v-model="profileForm.twitter" placeholder="Twitter 用户名 (选填)" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="cancelEdit">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import { Edit, Message, Calendar, Document, ChatLineRound, Link, Upload } from '@element-plus/icons-vue'

export default {
  name: 'UserProfile',
  components: {
    Edit,
    Message,
    Calendar,
    Document,
    ChatLineRound,
    Link,
    Upload
  },
  setup() {
    const store = useStore()
    const loading = ref(true)
    const isEditing = ref(false)
    const profileFormRef = ref(null)
    const fileInput = ref(null)
    const avatarUrl = ref('')
    const defaultAvatar = '/images/default-avatar.png'
    
    const user = computed(() => store.state.user.currentUser || {})
    
    const profileForm = reactive({
      username: '',
      email: '',
      bio: '',
      website: '',
      github: '',
      twitter: '',
      avatar: null
    })
    
    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
      ],
      website: [
        { pattern: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w.-]*)*\/?$/, message: '请输入有效的网址', trigger: 'blur' }
      ]
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`
    }
    
    // 获取用户资料
    const fetchUserProfile = async () => {
      loading.value = true
      try {
        await store.dispatch('getUserProfile')
        initForm()
      } catch (error) {
        console.error('获取用户资料失败', error)
        ElMessage.error('获取用户资料失败，请稍后再试')
      } finally {
        loading.value = false
      }
    }
    
    // 初始化表单
    const initForm = () => {
      profileForm.username = user.value.username || ''
      profileForm.email = user.value.email || ''
      profileForm.bio = user.value.bio || ''
      profileForm.website = user.value.website || ''
      profileForm.github = user.value.github || ''
      profileForm.twitter = user.value.twitter || ''
      avatarUrl.value = user.value.avatar || ''
    }
    
    // 触发文件上传
    const triggerUpload = () => {
      fileInput.value.click()
    }
    
    // 处理文件选择
    const handleFileChange = (e) => {
      const file = e.target.files[0]
      if (!file) return
      
      // 检查文件类型
      if (!file.type.match('image.*')) {
        ElMessage.error('请选择图片文件')
        return
      }
      
      // 预览图片
      const reader = new FileReader()
      reader.onload = (e) => {
        avatarUrl.value = e.target.result
      }
      reader.readAsDataURL(file)
      
      profileForm.avatar = file
    }
    
    // 提交表单
    const submitForm = async () => {
      if (!profileFormRef.value) return
      
      await profileFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        try {
          const formData = new FormData()
          
          formData.append('username', profileForm.username)
          formData.append('bio', profileForm.bio || '')
          formData.append('website', profileForm.website || '')
          formData.append('github', profileForm.github || '')
          formData.append('twitter', profileForm.twitter || '')
          
          if (profileForm.avatar) {
            formData.append('avatar', profileForm.avatar)
          }
          
          const result = await store.dispatch('updateUserProfile', formData)
          
          if (result.success) {
            ElMessage.success('资料更新成功')
            isEditing.value = false
            // 更新用户信息
            await store.dispatch('getUserProfile')
          } else {
            ElMessage.error(result.message || '更新失败')
          }
        } catch (error) {
          console.error('更新用户资料失败', error)
          ElMessage.error('更新用户资料失败，请稍后再试')
        }
      })
    }
    
    // 取消编辑
    const cancelEdit = () => {
      isEditing.value = false
      initForm()
    }
    
    onMounted(() => {
      fetchUserProfile()
    })
    
    return {
      loading,
      user,
      isEditing,
      profileForm,
      profileFormRef,
      fileInput,
      avatarUrl,
      defaultAvatar,
      rules,
      formatDate,
      triggerUpload,
      handleFileChange,
      submitForm,
      cancelEdit
    }
  }
}
</script>

<style scoped>
.user-profile {
  min-height: 80vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.loading-container {
  padding: 20px;
}

.profile-skeleton {
  display: flex;
  align-items: flex-start;
}

.avatar-skeleton {
  width: 100px;
  height: 100px;
  border-radius: 50%;
}

.skeleton-title {
  height: 24px;
  width: 40%;
}

.profile-display {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-header {
  display: flex;
  gap: 24px;
  align-items: center;
}

.profile-avatar {
  flex-shrink: 0;
}

.profile-avatar img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--el-color-primary-light-7);
}

.profile-basic {
  flex: 1;
}

.profile-name {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: var(--el-text-color-primary);
}

.profile-meta {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--el-text-color-secondary);
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-section {
  padding: 16px;
  background-color: var(--el-fill-color-blank);
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color-lighter);
  padding-bottom: 8px;
}

.bio-content {
  white-space: pre-line;
  line-height: 1.6;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.link-item a {
  color: var(--el-color-primary);
  text-decoration: none;
}

.link-item a:hover {
  text-decoration: underline;
}

.custom-icon {
  flex-shrink: 0;
}

.empty-links {
  color: var(--el-text-color-secondary);
  font-style: italic;
}

.profile-edit {
  max-width: 600px;
  margin: 0 auto;
}

.avatar-uploader {
  display: flex;
  gap: 16px;
  align-items: center;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .profile-meta {
    grid-template-columns: 1fr;
    justify-items: center;
  }
  
  .avatar-uploader {
    flex-direction: column;
  }
}
</style> 
<template>
  <div class="user-settings">
    <h2 class="page-title">账号设置</h2>
    
    <el-tabs type="border-card">
      <!-- 密码修改 -->
      <el-tab-pane>
        <template #label>
          <div class="tab-label">
            <el-icon><Lock /></el-icon>
            <span>密码修改</span>
          </div>
        </template>
        
        <div class="tab-content">
          <el-form
            :model="passwordForm"
            :rules="passwordRules"
            ref="passwordFormRef"
            label-width="120px"
            v-loading="passwordLoading"
          >
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="submitPasswordForm">
                修改密码
              </el-button>
              <el-button @click="resetPasswordForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      
      <!-- 主题设置 -->
      <el-tab-pane>
        <template #label>
          <div class="tab-label">
            <el-icon><Brush /></el-icon>
            <span>主题设置</span>
          </div>
        </template>
        
        <div class="tab-content">
          <div class="setting-section">
            <div class="setting-label">暗色模式</div>
            <div class="setting-control">
              <el-switch
                v-model="darkMode"
                @change="toggleTheme"
                active-text="开启"
                inactive-text="关闭"
              />
            </div>
          </div>
          
          <div class="setting-section">
            <div class="setting-label">主题颜色</div>
            <div class="setting-control">
              <div class="color-options">
                <div
                  v-for="color in themeColors"
                  :key="color.value"
                  class="color-option"
                  :class="{ active: themeColor === color.value }"
                  :style="{ backgroundColor: color.color }"
                  @click="setThemeColor(color.value)"
                ></div>
              </div>
            </div>
          </div>
          
          <div class="setting-section">
            <div class="setting-label">字体大小</div>
            <div class="setting-control">
              <el-slider
                v-model="fontSize"
                :min="12"
                :max="20"
                :step="1"
                show-stops
                @change="setFontSize"
              />
              <div class="font-size-labels">
                <span>小</span>
                <span>默认</span>
                <span>大</span>
              </div>
            </div>
          </div>
          
          <div class="setting-section">
            <el-button type="primary" @click="saveThemeSettings">
              保存设置
            </el-button>
            <el-button @click="resetThemeSettings">
              恢复默认
            </el-button>
          </div>
        </div>
      </el-tab-pane>
      
      <!-- 隐私设置 -->
      <el-tab-pane>
        <template #label>
          <div class="tab-label">
            <el-icon><Lock /></el-icon>
            <span>隐私设置</span>
          </div>
        </template>
        
        <div class="tab-content">
          <div class="setting-section">
            <div class="setting-label">公开个人资料</div>
            <div class="setting-control">
              <el-switch
                v-model="privacySettings.publicProfile"
                @change="updatePrivacySettings"
                active-text="允许"
                inactive-text="禁止"
              />
            </div>
          </div>
          
          <div class="setting-section">
            <div class="setting-label">显示我的评论</div>
            <div class="setting-control">
              <el-switch
                v-model="privacySettings.showComments"
                @change="updatePrivacySettings"
                active-text="允许"
                inactive-text="禁止"
              />
            </div>
          </div>
          
          <div class="setting-section">
            <div class="setting-label">接收邮件通知</div>
            <div class="setting-control">
              <el-switch
                v-model="privacySettings.emailNotification"
                @change="updatePrivacySettings"
                active-text="允许"
                inactive-text="禁止"
              />
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import { Lock, Brush } from '@element-plus/icons-vue'

export default {
  name: 'UserSettings',
  components: {
    Lock,
    Brush
  },
  setup() {
    const store = useStore()
    const passwordFormRef = ref(null)
    const passwordLoading = ref(false)
    
    // 密码表单
    const passwordForm = reactive({
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    // 密码表单验证规则
    const passwordRules = {
      oldPassword: [
        { required: true, message: '请输入当前密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ],
      newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认新密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== passwordForm.newPassword) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }
    
    // 主题设置
    const darkMode = computed({
      get: () => store.getters.isDarkMode,
      set: (value) => {}
    })
    
    const themeColor = ref('default')
    const fontSize = ref(14)
    
    const themeColors = [
      { value: 'default', color: '#409eff', label: '默认蓝' },
      { value: 'green', color: '#67c23a', label: '活力绿' },
      { value: 'red', color: '#f56c6c', label: '热情红' },
      { value: 'orange', color: '#e6a23c', label: '温暖橙' },
      { value: 'purple', color: '#9c27b0', label: '高贵紫' }
    ]
    
    // 隐私设置
    const privacySettings = reactive({
      publicProfile: true,
      showComments: true,
      emailNotification: true
    })
    
    // 切换主题
    const toggleTheme = () => {
      store.dispatch('toggleTheme')
    }
    
    // 设置主题颜色
    const setThemeColor = (color) => {
      themeColor.value = color
      // 实际应用主题颜色的逻辑
      document.documentElement.style.setProperty('--primary-color', getColorByValue(color))
    }
    
    // 获取颜色值
    const getColorByValue = (value) => {
      const colorObj = themeColors.find(c => c.value === value)
      return colorObj ? colorObj.color : '#409eff'
    }
    
    // 设置字体大小
    const setFontSize = (size) => {
      fontSize.value = size
      document.documentElement.style.fontSize = `${size}px`
    }
    
    // 保存主题设置
    const saveThemeSettings = () => {
      localStorage.setItem('themeColor', themeColor.value)
      localStorage.setItem('fontSize', fontSize.value)
      ElMessage.success('主题设置已保存')
    }
    
    // 重置主题设置
    const resetThemeSettings = () => {
      themeColor.value = 'default'
      fontSize.value = 14
      document.documentElement.style.setProperty('--primary-color', '#409eff')
      document.documentElement.style.fontSize = '14px'
      localStorage.removeItem('themeColor')
      localStorage.removeItem('fontSize')
      ElMessage.success('已恢复默认设置')
    }
    
    // 更新隐私设置
    const updatePrivacySettings = async () => {
      try {
        await store._vm.$axios.put('/user/privacy', privacySettings)
        ElMessage.success('隐私设置已更新')
      } catch (error) {
        console.error('更新隐私设置失败', error)
        ElMessage.error('更新隐私设置失败，请稍后再试')
      }
    }
    
    // 提交密码表单
    const submitPasswordForm = async () => {
      if (!passwordFormRef.value) return
      
      await passwordFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        passwordLoading.value = true
        
        try {
          const result = await store.dispatch('changePassword', {
            oldPassword: passwordForm.oldPassword,
            newPassword: passwordForm.newPassword
          })
          
          if (result.success) {
            ElMessage.success('密码修改成功')
            resetPasswordForm()
          } else {
            ElMessage.error(result.message || '密码修改失败，请稍后再试')
          }
        } catch (error) {
          console.error('修改密码失败', error)
          ElMessage.error('修改密码失败，请稍后再试')
        } finally {
          passwordLoading.value = false
        }
      })
    }
    
    // 重置密码表单
    const resetPasswordForm = () => {
      if (passwordFormRef.value) {
        passwordFormRef.value.resetFields()
      }
    }
    
    // 初始化设置
    const initSettings = () => {
      // 加载保存的主题颜色
      const savedThemeColor = localStorage.getItem('themeColor')
      if (savedThemeColor) {
        themeColor.value = savedThemeColor
        document.documentElement.style.setProperty('--primary-color', getColorByValue(savedThemeColor))
      }
      
      // 加载保存的字体大小
      const savedFontSize = localStorage.getItem('fontSize')
      if (savedFontSize) {
        fontSize.value = parseInt(savedFontSize, 10)
        document.documentElement.style.fontSize = `${fontSize.value}px`
      }
      
      // 加载隐私设置
      fetchPrivacySettings()
    }
    
    // 获取隐私设置
    const fetchPrivacySettings = async () => {
      try {
        const response = await store._vm.$axios.get('/user/privacy')
        Object.assign(privacySettings, response.data)
      } catch (error) {
        console.error('获取隐私设置失败', error)
      }
    }
    
    onMounted(() => {
      initSettings()
    })
    
    return {
      passwordFormRef,
      passwordForm,
      passwordRules,
      passwordLoading,
      darkMode,
      themeColor,
      fontSize,
      themeColors,
      privacySettings,
      toggleTheme,
      setThemeColor,
      setFontSize,
      saveThemeSettings,
      resetThemeSettings,
      updatePrivacySettings,
      submitPasswordForm,
      resetPasswordForm
    }
  }
}
</script>

<style scoped>
.user-settings {
  max-width: 800px;
  margin: 0 auto;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
  color: #333;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-content {
  padding: 20px 10px;
}

.setting-section {
  margin-bottom: 24px;
  display: flex;
  align-items: center;
}

.setting-label {
  width: 120px;
  font-weight: 500;
  color: #333;
}

.setting-control {
  flex: 1;
}

.color-options {
  display: flex;
  gap: 16px;
  margin-top: 10px;
}

.color-option {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  border: 2px solid transparent;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: #333;
}

.color-option.active::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: white;
}

.font-size-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}

@media (max-width: 576px) {
  .setting-section {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .setting-label {
    width: 100%;
    margin-bottom: 8px;
  }
  
  .setting-control {
    width: 100%;
  }
}
</style> 
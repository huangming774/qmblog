<template>
  <div class="register-page bg-gradient-to-br from-secondary/30 via-primary/20 to-accent/30">
    <div class="register-backdrop"></div>
    <div class="register-container">
      <div class="register-form backdrop-blur-md bg-white/60 shadow-2xl">
        <h2 class="form-title text-2xl font-bold text-gray-800 mb-8">注册</h2>
        
        <el-form 
          :model="registerForm" 
          :rules="rules" 
          ref="registerFormRef" 
          label-position="top"
          class="w-full"
        >
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="registerForm.username" 
              placeholder="请输入用户名"
              class="glass-input"
            />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="registerForm.email" 
              placeholder="请输入邮箱"
              class="glass-input"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="registerForm.password" 
              type="password" 
              placeholder="请输入密码"
              show-password
              class="glass-input"
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="registerForm.confirmPassword" 
              type="password" 
              placeholder="请再次输入密码"
              show-password
              class="glass-input"
            />
          </el-form-item>
          
          <div class="form-agreement mb-6">
            <el-checkbox v-model="agreement">
              我已阅读并同意 <el-link type="primary" :underline="false" class="hover:text-secondary">服务条款</el-link> 和 <el-link type="primary" :underline="false" class="hover:text-secondary">隐私政策</el-link>
            </el-checkbox>
          </div>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn w-full text-lg font-medium py-3 rounded-lg transition-all duration-300 hover:shadow-lg hover:opacity-90"
              :loading="loading"
              :disabled="!agreement"
              @click="handleRegister"
              v-motion="{ 
                initial: { scale: 0.98, opacity: 0.8 },
                hover: { scale: 1.02, opacity: 1 } 
              }"
            >
              注册
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="form-footer text-center mt-8 text-gray-600">
          <p>已有账号? <router-link to="/login" class="text-primary font-medium hover:text-secondary transition-colors">立即登录</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'Register',
  setup() {
    const store = useStore()
    const router = useRouter()
    
    const registerFormRef = ref(null)
    const loading = ref(false)
    const agreement = ref(false)
    
    // 注册表单
    const registerForm = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })
    
    // 密码验证一致性检查
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    
    // 表单验证规则
    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        { validator: validatePass, trigger: 'blur' }
      ]
    }
    
    // 处理注册
    const handleRegister = async () => {
      if (!registerFormRef.value) return
      
      await registerFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        loading.value = true
        
        try {
          const userData = {
            username: registerForm.username,
            email: registerForm.email,
            password: registerForm.password
          }
          
          const result = await store.dispatch('register', userData)
          
          if (result.success) {
            ElMessage({
              type: 'success',
              message: '注册成功！'
            })
            
            router.push('/')
          } else {
            ElMessage({
              type: 'error',
              message: result.message || '注册失败，请稍后再试'
            })
          }
        } catch (error) {
          console.error('注册出错', error)
          ElMessage({
            type: 'error',
            message: '注册失败，请稍后再试'
          })
        } finally {
          loading.value = false
        }
      })
    }
    
    return {
      registerFormRef,
      registerForm,
      rules,
      loading,
      agreement,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-page {
  @apply flex justify-center items-center min-h-[80vh] relative overflow-hidden py-6;
  background-size: 400% 400%;
  animation: gradientBG 15s ease infinite;
}

@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.register-backdrop {
  @apply absolute inset-0 z-0;
  background-image: url('data:image/svg+xml;charset=utf8,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"%3E%3Cpath fill="%236366f1" fill-opacity="0.15" d="M0,224L48,213.3C96,203,192,181,288,154.7C384,128,480,96,576,106.7C672,117,768,171,864,176C960,181,1056,139,1152,122.7C1248,107,1344,117,1392,122.7L1440,128L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"%3E%3C/path%3E%3C/svg%3E');
  background-size: cover;
  background-position: center bottom;
}

.register-container {
  @apply relative z-10 w-full max-w-md px-4;
}

.register-form {
  @apply p-8 rounded-2xl border border-white/30 transition-all duration-300;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1), 
              0 1px 8px rgba(0, 0, 0, 0.05),
              0 20px 30px -10px rgba(99, 102, 241, 0.2);
}

.form-title {
  background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-align: center;
}

:deep(.glass-input .el-input__inner) {
  @apply bg-white/50 backdrop-blur-sm border border-white/60 transition-all duration-300;
}

:deep(.glass-input .el-input__inner:focus) {
  @apply bg-white/70 border-primary/60 shadow-sm;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

:deep(.el-form-item__label) {
  @apply text-gray-700 font-medium;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  @apply bg-primary border-primary;
}

:deep(.el-checkbox__inner:hover) {
  @apply border-primary;
}

/* 动画效果 */
.register-form {
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
}

/* 毛玻璃颜色变量 */
:root {
  --color-primary: #3b82f6;
  --color-secondary: #6366f1;
}
</style> 
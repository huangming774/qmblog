<template>
  <div class="login-page bg-gradient-to-br from-secondary/30 via-primary/20 to-accent/30">
    <div class="login-backdrop"></div>
    <div class="login-container">
      <div class="login-form backdrop-blur-md bg-white/60 shadow-2xl">
        <h2 class="form-title text-2xl font-bold text-gray-800 mb-8">登录</h2>
        
        <el-form 
          :model="loginForm" 
          :rules="rules" 
          ref="loginFormRef" 
          label-position="top"
          class="w-full"
        >
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="loginForm.email" 
              placeholder="请输入邮箱"
              prefix-icon="el-icon-message"
              class="glass-input"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="请输入密码"
              prefix-icon="el-icon-lock"
              show-password
              class="glass-input"
            />
          </el-form-item>
          
          <div class="form-options flex justify-between items-center mb-6">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-link type="primary" :underline="false" class="hover:text-secondary">忘记密码?</el-link>
          </div>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn w-full text-lg font-medium py-3 rounded-lg transition-all duration-300 hover:shadow-lg hover:opacity-90"
              :loading="loading"
              @click="handleLogin"
              v-motion="{ 
                initial: { scale: 0.98, opacity: 0.8 },
                hover: { scale: 1.02, opacity: 1 } 
              }"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="form-footer text-center mt-8 text-gray-600">
          <p>还没有账号? <router-link to="/register" class="text-primary font-medium hover:text-secondary transition-colors">立即注册</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'Login',
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    
    const loginFormRef = ref(null)
    const loading = ref(false)
    const rememberMe = ref(false)
    
    // 登录表单
    const loginForm = reactive({
      email: '',
      password: ''
    })
    
    // 表单验证规则
    const rules = {
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ]
    }
    
    // 处理登录
    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      await loginFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        loading.value = true
        
        try {
          const result = await store.dispatch('login', loginForm)
          
          if (result.success) {
            ElMessage({
              type: 'success',
              message: '登录成功'
            })
            
            // 如果有重定向地址，跳转到重定向地址
            const redirectPath = route.query.redirect || '/'
            router.push(redirectPath)
          } else {
            ElMessage({
              type: 'error',
              message: result.message || '登录失败，请检查账号密码'
            })
          }
        } catch (error) {
          console.error('登录出错', error)
          ElMessage({
            type: 'error',
            message: '登录失败，请稍后再试'
          })
        } finally {
          loading.value = false
        }
      })
    }
    
    // 页面加载时如果已登录，跳转到首页
    onMounted(() => {
      if (store.state.isAuthenticated) {
        router.push('/')
      }
      
      // 演示账号（仅用于开发测试）
      loginForm.email = 'admin@example.com'
      loginForm.password = 'admin123'
    })
    
    return {
      loginFormRef,
      loginForm,
      rules,
      loading,
      rememberMe,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-page {
  @apply flex justify-center items-center min-h-[80vh] relative overflow-hidden;
  background-size: 400% 400%;
  animation: gradientBG 15s ease infinite;
}

@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.login-backdrop {
  @apply absolute inset-0 z-0;
  background-image: url('data:image/svg+xml;charset=utf8,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"%3E%3Cpath fill="%233b82f6" fill-opacity="0.15" d="M0,64L48,96C96,128,192,192,288,186.7C384,181,480,107,576,85.3C672,64,768,96,864,133.3C960,171,1056,213,1152,202.7C1248,192,1344,128,1392,96L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"%3E%3C/path%3E%3C/svg%3E');
  background-size: cover;
  background-position: center bottom;
}

.login-container {
  @apply relative z-10 w-full max-w-md px-4;
}

.login-form {
  @apply p-8 rounded-2xl border border-white/30 transition-all duration-300;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1), 
              0 1px 8px rgba(0, 0, 0, 0.05),
              0 20px 30px -10px rgba(59, 130, 246, 0.2);
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
.login-form {
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
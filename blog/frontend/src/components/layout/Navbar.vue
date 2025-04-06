<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    router
    class="navbar"
    :ellipsis="false"
  >
    <el-menu-item index="/" class="logo-container">
      <img src="@/assets/logo.png" alt="Logo" class="logo" v-if="false" />
      <span class="blog-title">博客系统</span>
    </el-menu-item>
    
    <div class="flex-grow" />
    
    <el-menu-item index="/">首页</el-menu-item>
    <el-menu-item index="/archive">归档</el-menu-item>
    <el-menu-item index="/tags">标签</el-menu-item>
    
    <!-- 搜索框 -->
    <div class="search-container">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索文章..."
        clearable
        @keyup.enter="handleSearch"
        class="search-input"
      >
        <template #suffix>
          <el-icon @click="handleSearch" class="search-icon">
            <Search />
          </el-icon>
        </template>
      </el-input>
    </div>
    
    <!-- 主题切换 -->
    <div class="theme-switch">
      <el-tooltip content="切换主题">
        <el-button 
          circle 
          :icon="isDarkMode ? 'Sunny' : 'Moon'" 
          @click="toggleTheme"
          class="theme-button"
        ></el-button>
      </el-tooltip>
    </div>
    
    <template v-if="isAuthenticated">
      <!-- 通知图标 -->
      <div class="notification-icon">
        <el-badge :is-dot="hasUnreadNotifications">
          <el-button 
            circle
            icon="Bell"
            @click="goToNotifications"
            class="notification-button"
          ></el-button>
        </el-badge>
      </div>
      
      <el-sub-menu index="user" v-if="currentUser">
        <template #title>
          <el-avatar :size="32" :src="currentUser.avatar || defaultAvatar">
            {{ currentUser.username ? currentUser.username.charAt(0).toUpperCase() : 'U' }}
          </el-avatar>
          <span class="username">{{ currentUser.username }}</span>
        </template>
        
        <el-menu-item index="/user/profile">
          <el-icon><User /></el-icon>
          个人资料
        </el-menu-item>
        
        <el-menu-item index="/user/posts">
          <el-icon><Document /></el-icon>
          我的文章
        </el-menu-item>
        
        <el-menu-item index="/user/favorites">
          <el-icon><Star /></el-icon>
          收藏列表
        </el-menu-item>
        
        <el-menu-item index="/user/settings">
          <el-icon><Setting /></el-icon>
          账号设置
        </el-menu-item>
        
        <el-menu-item index="/admin" v-if="isAdmin">
          <el-icon><Setting /></el-icon>
          管理后台
        </el-menu-item>
        
        <el-menu-item index="/admin/posts/create">
          <el-icon><EditPen /></el-icon>
          写文章
        </el-menu-item>
        
        <el-menu-item @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          退出登录
        </el-menu-item>
      </el-sub-menu>
    </template>
    
    <template v-else>
      <el-menu-item index="/login">登录</el-menu-item>
      <el-menu-item index="/register">注册</el-menu-item>
    </template>
  </el-menu>
</template>

<script>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { 
  Setting, 
  EditPen, 
  SwitchButton, 
  User, 
  Document, 
  Star, 
  Search, 
  Bell 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

export default {
  name: 'Navbar',
  components: {
    Setting,
    EditPen,
    SwitchButton,
    User,
    Document, 
    Star,
    Search,
    Bell
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const store = useStore()
    
    const activeIndex = ref('/')
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    const searchKeyword = ref('')
    
    // 从vuex获取状态
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    const currentUser = computed(() => store.state.user)
    const isAdmin = computed(() => store.getters.isAdmin)
    const isDarkMode = computed(() => store.getters.isDarkMode)
    const hasUnreadNotifications = computed(() => store.getters.hasUnreadNotifications)
    
    // 监听路由变化更新activeIndex
    const updateActiveIndex = () => {
      const path = route.path
      if (path === '/') {
        activeIndex.value = '/'
      } else if (path.startsWith('/admin')) {
        activeIndex.value = '/admin'
      } else {
        activeIndex.value = path
      }
    }
    
    // 初始化
    updateActiveIndex()
    
    // 搜索处理
    const handleSearch = () => {
      if (searchKeyword.value.trim()) {
        router.push({
          path: '/search',
          query: { keyword: searchKeyword.value }
        })
      }
    }
    
    // 切换主题
    const toggleTheme = () => {
      store.dispatch('toggleTheme')
    }
    
    // 前往通知页面
    const goToNotifications = () => {
      router.push('/user/notifications')
    }
    
    // 退出登录
    const handleLogout = () => {
      store.dispatch('logout')
      router.push('/')
      ElMessage({
        message: '您已成功退出登录',
        type: 'success'
      })
    }
    
    return {
      activeIndex,
      isAuthenticated,
      currentUser,
      isAdmin,
      defaultAvatar,
      searchKeyword,
      isDarkMode,
      hasUnreadNotifications,
      handleSearch,
      toggleTheme,
      goToNotifications,
      handleLogout
    }
  }
}
</script>

<style scoped>
.navbar {
  padding: 0 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.logo-container {
  display: flex;
  align-items: center;
}

.logo {
  height: 40px;
  margin-right: 10px;
}

.blog-title {
  font-size: 20px;
  font-weight: bold;
  color: var(--primary-color);
}

.flex-grow {
  flex-grow: 1;
}

.search-container {
  margin: 0 15px;
  display: flex;
  align-items: center;
}

.search-input {
  width: 200px;
  transition: width 0.3s;
}

.search-input:focus-within {
  width: 250px;
}

.search-icon {
  cursor: pointer;
}

.theme-switch, .notification-icon {
  margin: 0 10px;
  display: flex;
  align-items: center;
}

.theme-button, .notification-button {
  font-size: 18px;
}

.username {
  margin-left: 8px;
  font-size: 14px;
}

@media (max-width: 992px) {
  .search-input {
    width: 150px;
  }
  
  .search-input:focus-within {
    width: 200px;
  }
}

@media (max-width: 768px) {
  .username {
    display: none;
  }
  
  .search-input {
    width: 120px;
  }
}

@media (max-width: 576px) {
  .search-container {
    display: none;
  }
}
</style> 
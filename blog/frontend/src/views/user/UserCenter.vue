<template>
  <div class="user-center">
    <el-container class="container">
      <el-aside width="250px" class="aside">
        <div class="user-info">
          <el-avatar 
            :size="80" 
            :src="currentUser.avatar || defaultAvatar"
            class="avatar"
          >
            {{ currentUser.username ? currentUser.username.charAt(0).toUpperCase() : 'U' }}
          </el-avatar>
          <h3 class="username">{{ currentUser.username }}</h3>
          <p class="role">{{ currentUser.role === 'admin' ? '管理员' : '普通用户' }}</p>
        </div>
        
        <el-menu
          :default-active="activeMenu"
          router
          class="side-menu"
        >
          <el-menu-item index="/user/profile">
            <el-icon><User /></el-icon>
            <span>个人资料</span>
          </el-menu-item>
          
          <el-menu-item index="/user/posts">
            <el-icon><Document /></el-icon>
            <span>我的文章</span>
          </el-menu-item>
          
          <el-menu-item index="/user/favorites">
            <el-icon><Star /></el-icon>
            <span>收藏列表</span>
          </el-menu-item>
          
          <el-menu-item index="/user/notifications">
            <el-icon><Bell /></el-icon>
            <span>
              消息通知
              <el-badge 
                :value="unreadNotificationsCount" 
                :hidden="unreadNotificationsCount === 0" 
                class="notification-badge"
              />
            </span>
          </el-menu-item>
          
          <el-menu-item index="/user/settings">
            <el-icon><Setting /></el-icon>
            <span>账号设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { User, Document, Star, Bell, Setting } from '@element-plus/icons-vue'

export default {
  name: 'UserCenter',
  components: {
    User,
    Document,
    Star,
    Bell,
    Setting
  },
  setup() {
    const route = useRoute()
    const store = useStore()
    
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    
    // 从vuex获取状态
    const currentUser = computed(() => store.state.currentUser || {})
    const unreadNotificationsCount = computed(() => store.state.unreadNotificationsCount)
    
    // 激活的菜单项
    const activeMenu = computed(() => route.path)
    
    onMounted(() => {
      // 获取通知
      store.dispatch('fetchNotifications')
    })
    
    return {
      currentUser,
      defaultAvatar,
      activeMenu,
      unreadNotificationsCount
    }
  }
}
</script>

<style scoped>
.user-center {
  min-height: calc(100vh - 60px);
  padding: 20px;
  background-color: #f5f7fa;
}

.container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.aside {
  background-color: #fff;
  border-right: 1px solid #ebeef5;
}

.user-info {
  padding: 30px 20px;
  text-align: center;
  border-bottom: 1px solid #ebeef5;
}

.avatar {
  margin-bottom: 15px;
  border: 3px solid rgba(var(--primary-color-rgb), 0.2);
}

.username {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 5px;
  color: #333;
}

.role {
  font-size: 14px;
  color: #909399;
}

.side-menu {
  border-right: none;
}

.main {
  padding: 20px;
  background-color: #fff;
}

.notification-badge {
  margin-left: 8px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .user-center {
    padding: 10px;
  }
  
  .aside {
    width: 100% !important;
  }
  
  .container {
    flex-direction: column;
  }
}
</style> 
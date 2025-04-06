<template>
  <div class="user-sidebar">
    <div class="user-info">
      <div class="avatar-container">
        <img :src="user.avatar || defaultAvatar" :alt="user.username" class="user-avatar" />
      </div>
      <h3 class="username">{{ user.username }}</h3>
      <p class="bio">{{ user.bio || '这个人很懒，什么都没写...' }}</p>
    </div>
    
    <div class="menu-container">
      <el-menu
        :default-active="activeMenu"
        class="user-menu"
        router
      >
        <el-menu-item index="/user/profile">
          <el-icon><User /></el-icon>
          <span>个人资料</span>
        </el-menu-item>
        
        <el-menu-item index="/user/posts">
          <el-icon><Document /></el-icon>
          <span>我的文章</span>
        </el-menu-item>
        
        <el-menu-item index="/user/comments">
          <el-icon><ChatLineRound /></el-icon>
          <span>我的评论</span>
        </el-menu-item>
        
        <el-menu-item index="/user/favorites">
          <el-icon><Star /></el-icon>
          <span>我的收藏</span>
        </el-menu-item>
        
        <el-menu-item index="/user/settings">
          <el-icon><Setting /></el-icon>
          <span>账号设置</span>
        </el-menu-item>
      </el-menu>
    </div>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import { User, Document, ChatLineRound, Star, Setting } from '@element-plus/icons-vue'

export default {
  name: 'UserSidebar',
  components: {
    User,
    Document,
    ChatLineRound,
    Star,
    Setting
  },
  setup() {
    const store = useStore()
    const route = useRoute()
    
    const user = computed(() => store.state.user.currentUser || {})
    const defaultAvatar = '/images/default-avatar.png'
    
    const activeMenu = computed(() => {
      return route.path
    })
    
    return {
      user,
      defaultAvatar,
      activeMenu
    }
  }
}
</script>

<style scoped>
.user-sidebar {
  border-radius: 8px;
  background-color: var(--el-bg-color);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.user-info {
  padding: 24px;
  text-align: center;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.avatar-container {
  margin-bottom: 16px;
}

.user-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--el-color-primary-light-7);
}

.username {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.bio {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.menu-container {
  padding: 16px 0;
}

.user-menu {
  border-right: none;
}

@media (max-width: 768px) {
  .user-sidebar {
    margin-bottom: 20px;
  }
}
</style> 
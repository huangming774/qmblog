<template>
  <div class="user-notifications">
    <div class="page-header">
      <h2 class="page-title">消息通知</h2>
      
      <div class="header-actions">
        <el-button 
          type="primary" 
          plain 
          @click="markAllAsRead" 
          :disabled="!hasUnreadNotifications"
        >
          <el-icon><Check /></el-icon>
          全部标为已读
        </el-button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>
    
    <div v-else-if="notifications.length === 0" class="empty-container">
      <el-empty description="暂无消息通知" :image-size="200">
        <template #description>
          <p>您的消息中心暂时没有通知</p>
        </template>
      </el-empty>
    </div>
    
    <div v-else class="notifications-list">
      <!-- 未读通知 -->
      <div v-if="unreadNotifications.length > 0" class="notification-section">
        <div class="section-header">
          <h3 class="section-title">未读通知</h3>
          <el-tag type="danger" class="count-tag">{{ unreadNotifications.length }}</el-tag>
        </div>
        
        <div class="notification-group">
          <el-card
            v-for="notification in unreadNotifications"
            :key="notification.id"
            class="notification-card unread"
            shadow="hover"
          >
            <div class="notification-content">
              <div class="notification-icon">
                <el-avatar :size="40" :icon="getNotificationIcon(notification.type)" />
              </div>
              
              <div class="notification-details">
                <div class="notification-header">
                  <div class="notification-title">{{ getNotificationTitle(notification) }}</div>
                  <div class="notification-time">{{ formatTime(notification.createdAt) }}</div>
                </div>
                
                <div class="notification-message">{{ notification.content }}</div>
                
                <div class="notification-actions">
                  <el-button 
                    type="primary" 
                    link 
                    @click="handleNotificationClick(notification)"
                  >
                    查看详情
                  </el-button>
                  
                  <el-button 
                    type="info" 
                    link 
                    @click="markAsRead(notification.id)"
                  >
                    标为已读
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      
      <!-- 已读通知 -->
      <div v-if="readNotifications.length > 0" class="notification-section">
        <div class="section-header">
          <h3 class="section-title">已读通知</h3>
          <el-tag type="info" class="count-tag">{{ readNotifications.length }}</el-tag>
        </div>
        
        <div class="notification-group">
          <el-card
            v-for="notification in readNotifications"
            :key="notification.id"
            class="notification-card read"
            shadow="hover"
          >
            <div class="notification-content">
              <div class="notification-icon">
                <el-avatar :size="40" :icon="getNotificationIcon(notification.type)" />
              </div>
              
              <div class="notification-details">
                <div class="notification-header">
                  <div class="notification-title">{{ getNotificationTitle(notification) }}</div>
                  <div class="notification-time">{{ formatTime(notification.createdAt) }}</div>
                </div>
                
                <div class="notification-message">{{ notification.content }}</div>
                
                <div class="notification-actions">
                  <el-button 
                    type="primary" 
                    link 
                    @click="handleNotificationClick(notification)"
                  >
                    查看详情
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      
      <div class="pagination-container" v-if="total > pageSize">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          :current-page="currentPage"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Check, ChatLineRound, ThumbUp, Star, Bell } from '@element-plus/icons-vue'

export default {
  name: 'UserNotifications',
  components: {
    Check,
    ChatLineRound,
    ThumbUp,
    Star,
    Bell
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const loading = ref(false)
    const notifications = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    
    // 已读和未读通知
    const unreadNotifications = computed(() => 
      notifications.value.filter(n => !n.read)
    )
    
    const readNotifications = computed(() => 
      notifications.value.filter(n => n.read)
    )
    
    const hasUnreadNotifications = computed(() => 
      unreadNotifications.value.length > 0
    )
    
    // 获取通知图标
    const getNotificationIcon = (type) => {
      switch (type) {
        case 'comment':
          return ChatLineRound
        case 'like':
          return ThumbUp
        case 'favorite':
          return Star
        default:
          return Bell
      }
    }
    
    // 获取通知标题
    const getNotificationTitle = (notification) => {
      switch (notification.type) {
        case 'comment':
          return '评论通知'
        case 'like':
          return '点赞通知'
        case 'favorite':
          return '收藏通知'
        case 'system':
          return '系统通知'
        default:
          return '新通知'
      }
    }
    
    // 格式化时间
    const formatTime = (timeString) => {
      if (!timeString) return ''
      
      const now = new Date()
      const date = new Date(timeString)
      const diff = (now - date) / 1000 // 差值，单位为秒
      
      if (diff < 60) {
        return '刚刚'
      } else if (diff < 3600) {
        return `${Math.floor(diff / 60)}分钟前`
      } else if (diff < 86400) {
        return `${Math.floor(diff / 3600)}小时前`
      } else if (diff < 2592000) {
        return `${Math.floor(diff / 86400)}天前`
      } else {
        return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`
      }
    }
    
    // 获取通知列表
    const fetchNotifications = async () => {
      loading.value = true
      try {
        const response = await store._vm.$axios.get('/user/notifications', {
          params: {
            page: currentPage.value,
            pageSize: pageSize.value
          }
        })
        
        notifications.value = response.data.data || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('获取通知列表失败', error)
        ElMessage.error('获取通知列表失败，请稍后再试')
      } finally {
        loading.value = false
      }
    }
    
    // 标记通知为已读
    const markAsRead = async (notificationId) => {
      try {
        const result = await store.dispatch('markNotificationRead', notificationId)
        
        if (result.success) {
          // 更新本地状态
          const notification = notifications.value.find(n => n.id === notificationId)
          if (notification) {
            notification.read = true
          }
        } else {
          ElMessage.error(result.message || '标记失败')
        }
      } catch (error) {
        console.error('标记通知失败', error)
        ElMessage.error('标记通知失败，请稍后再试')
      }
    }
    
    // 标记所有通知为已读
    const markAllAsRead = async () => {
      try {
        const result = await store.dispatch('markAllNotificationsRead')
        
        if (result.success) {
          // 更新本地状态
          notifications.value.forEach(notification => {
            notification.read = true
          })
          ElMessage.success('已全部标为已读')
        } else {
          ElMessage.error(result.message || '操作失败')
        }
      } catch (error) {
        console.error('标记所有通知失败', error)
        ElMessage.error('操作失败，请稍后再试')
      }
    }
    
    // 处理通知点击
    const handleNotificationClick = (notification) => {
      // 标记为已读
      if (!notification.read) {
        markAsRead(notification.id)
      }
      
      // 根据通知类型和关联数据跳转
      if (notification.data) {
        try {
          const data = typeof notification.data === 'string' 
            ? JSON.parse(notification.data) 
            : notification.data
          
          if (notification.type === 'comment' && data.postId) {
            router.push(`/posts/${data.postId}#comment-${data.commentId}`)
          } else if ((notification.type === 'like' || notification.type === 'favorite') && data.postId) {
            router.push(`/posts/${data.postId}`)
          } else if (notification.type === 'system' && data.url) {
            router.push(data.url)
          }
        } catch (error) {
          console.error('解析通知数据失败', error)
        }
      }
    }
    
    // 切换页码
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchNotifications()
    }
    
    onMounted(() => {
      fetchNotifications()
    })
    
    return {
      loading,
      notifications,
      unreadNotifications,
      readNotifications,
      hasUnreadNotifications,
      currentPage,
      pageSize,
      total,
      getNotificationIcon,
      getNotificationTitle,
      formatTime,
      markAsRead,
      markAllAsRead,
      handleNotificationClick,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.user-notifications {
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

.header-actions {
  display: flex;
  gap: 10px;
}

.loading-container,
.empty-container {
  padding: 40px 0;
}

.notification-section {
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  margin-right: 10px;
}

.count-tag {
  font-size: 12px;
}

.notification-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-card {
  transition: transform 0.2s;
}

.notification-card:hover {
  transform: translateY(-2px);
}

.notification-card.unread {
  border-left: 4px solid var(--el-color-danger);
}

.notification-card.read {
  opacity: 0.8;
}

.notification-content {
  display: flex;
  gap: 16px;
}

.notification-icon {
  flex-shrink: 0;
}

.notification-details {
  flex: 1;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.notification-title {
  font-weight: 600;
  font-size: 16px;
}

.notification-time {
  color: #909399;
  font-size: 13px;
}

.notification-message {
  color: #606266;
  margin-bottom: 12px;
  line-height: 1.5;
}

.notification-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 576px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .notification-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
</style>

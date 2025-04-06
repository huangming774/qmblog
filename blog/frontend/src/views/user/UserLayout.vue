<template>
  <div class="user-layout">
    <div class="container">
      <div class="row">
        <div class="sidebar-col">
          <UserSidebar />
        </div>
        <div class="content-col">
          <div class="content-container">
            <router-view />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import UserSidebar from '@/components/user/UserSidebar.vue'

export default {
  name: 'UserLayout',
  components: {
    UserSidebar
  }
}
</script>

<style scoped>
.user-layout {
  padding: 24px 0;
  min-height: calc(100vh - 60px - 200px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
}

.row {
  display: flex;
  gap: 24px;
}

.sidebar-col {
  flex: 0 0 250px;
}

.content-col {
  flex: 1;
}

.content-container {
  background-color: var(--el-bg-color);
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  padding: 24px;
}

@media (max-width: 768px) {
  .row {
    flex-direction: column;
  }
  
  .sidebar-col {
    flex: auto;
  }
}
</style> 
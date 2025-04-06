<template>
  <div class="post-card group">
    <div class="post-cover" v-if="post.cover">
      <img 
        :src="post.cover" 
        :alt="post.title" 
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-70"></div>
    </div>
    
    <div class="post-content">
      <router-link :to="`/posts/${post.id}`" class="post-title">
        {{ post.title }}
      </router-link>
      
      <div class="post-meta">
        <div class="post-meta-item">
          <el-icon><User /></el-icon>
          <span>{{ post.user ? post.user.username : '未知作者' }}</span>
        </div>
        
        <div class="post-meta-item">
          <el-icon><Calendar /></el-icon>
          <span>{{ formatDate(post.createdAt) }}</span>
        </div>
        
        <div class="post-meta-item">
          <el-icon><View /></el-icon>
          <span>{{ post.viewCount }} 阅读</span>
        </div>
      </div>
      
      <div class="post-summary">
        {{ post.summary || truncateContent(post.content) }}
      </div>
      
      <div class="post-tags" v-if="post.tags && post.tags.length > 0">
        <el-tag 
          v-for="tag in post.tags" 
          :key="tag.id" 
          size="small" 
          effect="plain"
          class="tag"
        >
          {{ tag.name }}
        </el-tag>
      </div>
      
      <div class="post-actions">
        <el-button 
          type="primary" 
          text 
          class="transition-all duration-300 group-hover:translate-x-1"
          @click="$router.push(`/posts/${post.id}`)"
        >
          阅读全文
          <el-icon class="el-icon--right"><ArrowRight /></el-icon>
        </el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { ArrowRight, Calendar, User, View } from '@element-plus/icons-vue'

export default {
  name: 'PostCard',
  components: {
    ArrowRight,
    Calendar,
    User,
    View
  },
  props: {
    post: {
      type: Object,
      required: true
    }
  },
  methods: {
    formatDate(dateString) {
      return this.$moment(dateString).format('YYYY-MM-DD')
    },
    truncateContent(content) {
      if (!content) return ''
      
      // 去除HTML标签
      const plainText = content.replace(/<[^>]+>/g, '')
      
      // 截取前150个字符
      return plainText.length > 150 
        ? plainText.substring(0, 150) + '...' 
        : plainText
    }
  }
}
</script>

<style scoped>
.post-card {
  @apply bg-white rounded-xl shadow-md overflow-hidden transition-all duration-300 hover:shadow-xl h-full flex flex-col;
}

.post-cover {
  @apply h-48 overflow-hidden relative;
}

.post-content {
  @apply p-5 flex flex-col flex-grow;
}

.post-title {
  @apply text-xl font-semibold text-gray-800 mb-3 no-underline transition-colors duration-200 hover:text-primary overflow-hidden line-clamp-2;
}

.post-meta {
  @apply flex text-gray-500 text-sm mb-4 flex-wrap;
}

.post-meta-item {
  @apply mr-4 flex items-center mb-1;
}

.post-meta-item .el-icon {
  @apply mr-1 text-xs;
}

.post-summary {
  @apply text-gray-600 text-base leading-relaxed mb-4 overflow-hidden flex-grow line-clamp-3;
}

.post-tags {
  @apply mb-4 flex flex-wrap;
}

.tag {
  @apply mr-2 mb-2 transition-transform duration-200 hover:scale-105;
}

.post-actions {
  @apply flex justify-end mt-auto pt-2 border-t border-gray-100;
}

@media (max-width: 768px) {
  .post-cover {
    @apply h-36;
  }
  
  .post-content {
    @apply p-4;
  }
  
  .post-title {
    @apply text-lg mb-2;
  }
  
  .post-summary {
    @apply text-sm line-clamp-2 mb-3;
  }
  
  .post-meta {
    @apply mb-3;
  }
}
</style> 
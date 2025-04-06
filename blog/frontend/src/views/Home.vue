<template>
  <div class="home">
    <div class="banner bg-gradient-to-r from-primary/80 to-secondary/70">
      <div class="banner-content" data-aos="fade-up">
        <h1 class="banner-title text-4xl md:text-5xl text-white font-bold mb-4">欢迎访问博客系统</h1>
        <p class="banner-description text-xl text-white/90 mb-8">探索知识的海洋，分享你的见解与创意</p>
        <el-button 
          type="primary" 
          size="large" 
          class="btn-primary text-lg px-8 py-3 shadow-lg"
          @click="isAuthenticated ? $router.push('/admin/posts/create') : $router.push('/login')"
          v-motion="{ initial: { y: 20, opacity: 0 }, enter: { y: 0, opacity: 1 } }"
        >
          {{ isAuthenticated ? '写文章' : '立即登录' }}
        </el-button>
      </div>
    </div>
    
    <div class="content-container max-w-7xl mx-auto px-4">
      <el-row :gutter="20">
        <el-col :span="18">
          <!-- 文章列表 -->
          <div class="posts-container" data-aos="fade-up" data-aos-delay="100">
            <div class="posts-header flex items-center justify-between mb-6">
              <h2 class="section-title text-2xl font-bold">最新文章</h2>
              
              <div class="posts-filter">
                <el-select v-model="selectedTag" placeholder="按标签筛选" clearable>
                  <el-option
                    v-for="tag in availableTags"
                    :key="tag.id"
                    :label="tag.name"
                    :value="tag.name"
                  />
                </el-select>
              </div>
            </div>
            
            <el-empty 
              description="暂无文章" 
              v-if="posts.length === 0 && !loading"
            />
            
            <el-skeleton :rows="3" animated v-if="loading" />
            
            <el-row :gutter="20" v-else>
              <el-col 
                :xs="24" 
                :sm="24" 
                :md="12" 
                :lg="8" 
                v-for="(post, index) in posts" 
                :key="post.id"
              >
                <div 
                  data-aos="fade-up" 
                  :data-aos-delay="index * 100"
                  v-motion="{ 
                    initial: { scale: 0.95, opacity: 0 }, 
                    enter: { 
                      scale: 1, 
                      opacity: 1,
                      transition: { delay: index * 50 }
                    },
                    hovered: { scale: 1.02 }
                  }"
                >
                  <PostCard :post="post" class="card" />
                </div>
              </el-col>
            </el-row>
            
            <div class="pagination-container mt-10 flex justify-center" v-if="totalPosts > 0">
              <el-pagination
                background
                layout="prev, pager, next"
                :total="totalPosts"
                :page-size="pageSize"
                :current-page="currentPage"
                @current-change="handlePageChange"
              />
            </div>
          </div>
        </el-col>
        
        <el-col :span="6">
          <!-- 侧边栏 -->
          <div class="sidebar">
            <!-- 关于博主 -->
            <div class="sidebar-section about-section card p-5 mb-6" data-aos="fade-left">
              <h3 class="sidebar-title text-xl font-bold mb-4">关于博主</h3>
              <div class="about-content">
                <img 
                  src="https://via.placeholder.com/100" 
                  alt="博主头像" 
                  class="about-avatar rounded-full mx-auto mb-4 border-4 border-light shadow-md"
                />
                <p class="about-description text-center text-gray-600">
                  欢迎来到我的博客！这里记录了我的学习和工作经验，希望能对你有所帮助。
                </p>
              </div>
            </div>
            
            <!-- 标签云 -->
            <div class="sidebar-section card p-5 mb-6" data-aos="fade-left" data-aos-delay="100">
              <h3 class="sidebar-title text-xl font-bold mb-4">标签云</h3>
              <div class="tags-cloud flex flex-wrap gap-2">
                <el-tag
                  v-for="(tag, index) in availableTags"
                  :key="tag.id"
                  :type="getRandomTagType()"
                  effect="plain"
                  class="tag-item cursor-pointer mb-2 transition-all"
                  @click="selectedTag = tag.name"
                  v-motion="{ 
                    initial: { scale: 0.8, opacity: 0 }, 
                    enter: { 
                      scale: 1, 
                      opacity: 1,
                      transition: { delay: index * 50 }
                    },
                    hovered: { scale: 1.1 }
                  }"
                >
                  {{ tag.name }}
                </el-tag>
              </div>
            </div>
            
            <!-- 热门文章 -->
            <div class="sidebar-section card p-5" data-aos="fade-left" data-aos-delay="200">
              <h3 class="sidebar-title text-xl font-bold mb-4">热门文章</h3>
              <ul class="popular-posts">
                <li 
                  v-for="(post, index) in popularPosts" 
                  :key="post.id"
                  class="popular-post flex justify-between items-center py-3 border-b border-gray-100 cursor-pointer hover:bg-gray-50 rounded px-2 transition-colors"
                  @click="$router.push(`/posts/${post.id}`)"
                  v-motion="{ 
                    initial: { x: -20, opacity: 0 }, 
                    enter: { 
                      x: 0, 
                      opacity: 1,
                      transition: { delay: index * 100 }
                    }
                  }"
                >
                  <div class="popular-post-title truncate mr-2">{{ post.title }}</div>
                  <div class="popular-post-meta flex items-center text-gray-500 text-sm whitespace-nowrap">
                    <i class="el-icon-view mr-1"></i>
                    <span>{{ post.viewCount }}</span>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useStore } from 'vuex'
import { useMotion } from '@vueuse/motion'
import PostCard from '@/components/posts/PostCard.vue'

export default {
  name: 'Home',
  components: {
    PostCard
  },
  setup() {
    const store = useStore()
    const loading = ref(true)
    const currentPage = ref(1)
    const pageSize = ref(9)
    const selectedTag = ref('')
    
    // 从store获取数据
    const posts = computed(() => store.state.posts)
    const totalPosts = computed(() => store.state.totalPosts)
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    
    // 所有可用标签
    const availableTags = ref([
      { id: 1, name: 'Go' },
      { id: 2, name: 'Vue' },
      { id: 3, name: 'React' },
      { id: 4, name: 'PostgreSQL' },
      { id: 5, name: 'Redis' }
    ])
    
    // 热门文章
    const popularPosts = ref([])
    
    // 加载文章
    const loadPosts = async () => {
      loading.value = true
      
      await store.dispatch('fetchPosts', {
        page: currentPage.value,
        pageSize: pageSize.value,
        status: 'published',
        tag: selectedTag.value
      })
      
      loading.value = false
    }
    
    // 获取热门文章
    const getPopularPosts = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/v1/posts?pageSize=5&sort=viewCount&order=desc')
        const data = await response.json()
        popularPosts.value = data.data
      } catch (error) {
        console.error('获取热门文章失败', error)
      }
    }
    
    // 分页变化处理
    const handlePageChange = (page) => {
      currentPage.value = page
      loadPosts()
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
    
    // 随机获取标签类型
    const getRandomTagType = () => {
      const types = ['', 'success', 'info', 'warning', 'danger']
      return types[Math.floor(Math.random() * types.length)]
    }
    
    // 监听标签变化
    watch(selectedTag, () => {
      currentPage.value = 1
      loadPosts()
    })
    
    // 页面加载时获取数据
    onMounted(() => {
      loadPosts()
      getPopularPosts()
    })
    
    return {
      posts,
      totalPosts,
      loading,
      currentPage,
      pageSize,
      isAuthenticated,
      selectedTag,
      availableTags,
      popularPosts,
      handlePageChange,
      getRandomTagType
    }
  }
}
</script>

<style scoped>
.banner {
  @apply min-h-[300px] flex items-center justify-center mb-10 rounded-xl text-center relative overflow-hidden shadow-xl;
}

.banner-content {
  @apply z-10 px-4 py-8;
}

.content-container {
  @apply pb-16;
}

.sidebar-section {
  @apply transition-all duration-300 hover:shadow-lg;
}

.about-avatar {
  @apply transition-all duration-300 hover:scale-105;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .banner {
    @apply min-h-[200px] mb-6;
  }
  
  .banner-title {
    @apply text-3xl;
  }
  
  .banner-description {
    @apply text-base;
  }
}
</style> 
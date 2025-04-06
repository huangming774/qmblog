<template>
  <div class="post-detail">
    <div v-if="loading">
      <el-skeleton animated :rows="15" />
    </div>
    
    <div v-else-if="!post">
      <el-empty description="文章不存在或已被删除" />
    </div>
    
    <div v-else class="article-container">
      <!-- 文章封面 -->
      <img 
        v-if="post.cover" 
        :src="post.cover" 
        :alt="post.title" 
        class="article-cover" 
      />
      
      <!-- 文章标题 -->
      <h1 class="article-title">{{ post.title }}</h1>
      
      <!-- 文章元信息 -->
      <div class="article-meta">
        <div class="article-meta-item">
          <i class="far fa-user"></i>
          <span>{{ post.user ? post.user.username : '未知作者' }}</span>
        </div>
        
        <div class="article-meta-item">
          <i class="far fa-clock"></i>
          <span>{{ formatDate(post.createdAt) }}</span>
        </div>
        
        <div class="article-meta-item">
          <i class="far fa-eye"></i>
          <span>{{ post.viewCount }} 阅读</span>
        </div>
        
        <div class="article-meta-item">
          <i class="far fa-comment"></i>
          <span>{{ post.comments ? post.comments.length : 0 }} 评论</span>
        </div>
      </div>
      
      <!-- 文章标签 -->
      <div class="article-tags" v-if="post.tags && post.tags.length > 0">
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
      
      <!-- 文章内容 -->
      <div class="article-content" v-html="renderedContent"></div>
      
      <!-- 文章底部 -->
      <div class="article-footer">
        <el-divider>
          <i class="fas fa-pencil-alt"></i> 文章结束
        </el-divider>
        
        <div class="article-actions">
          <el-button type="primary" @click="handleLike" :disabled="liked">
            <i class="far" :class="liked ? 'fas fa-thumbs-up' : 'far fa-thumbs-up'"></i>
            {{ liked ? '已赞' : '点赞' }}
          </el-button>
          
          <el-button type="success" @click="handleShare">
            <i class="fas fa-share-alt"></i> 分享
          </el-button>
          
          <el-button 
            v-if="canEdit" 
            type="warning" 
            @click="$router.push(`/admin/posts/edit/${post.id}`)"
          >
            <i class="fas fa-edit"></i> 编辑
          </el-button>
        </div>
      </div>
      
      <!-- 评论区 -->
      <CommentSection 
        :post-id="post.id" 
        :comments="post.comments || []" 
        @refresh-comments="fetchPost"
      />
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { marked } from 'marked'
import CommentSection from '@/components/posts/CommentSection.vue'
import highlightjs from 'highlight.js'
import 'highlight.js/styles/github.css'

export default {
  name: 'PostDetail',
  components: {
    CommentSection
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const store = useStore()
    
    const postId = computed(() => route.params.id)
    const loading = ref(true)
    const post = ref(null)
    const liked = ref(false)
    
    // 从vuex获取状态
    const currentUser = computed(() => store.state.user)
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    
    // 判断当前用户是否可以编辑文章
    const canEdit = computed(() => {
      if (!isAuthenticated.value || !currentUser.value || !post.value) return false
      
      return currentUser.value.role === 'admin' || currentUser.value.id === post.value.userId
    })
    
    // 渲染Markdown内容
    const renderedContent = computed(() => {
      if (!post.value || !post.value.content) return ''
      
      marked.setOptions({
        highlight: function(code, lang) {
          const language = highlightjs.getLanguage(lang) ? lang : 'plaintext';
          return highlightjs.highlight(code, { language }).value;
        },
        langPrefix: 'hljs language-',
        gfm: true,
        breaks: true
      })
      
      return marked(post.value.content)
    })
    
    // 格式化日期
    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }
    
    // 获取文章详情
    const fetchPost = async () => {
      loading.value = true
      
      try {
        const result = await store.dispatch('fetchPost', postId.value)
        
        if (result.success) {
          post.value = store.state.currentPost
          
          // 检查本地存储中是否已点赞
          const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
          liked.value = likedPosts.includes(parseInt(postId.value))
        } else {
          router.push('/404')
        }
      } catch (error) {
        console.error('获取文章详情失败', error)
        ElMessage.error('获取文章详情失败，请稍后再试')
      } finally {
        loading.value = false
      }
    }
    
    // 点赞处理
    const handleLike = () => {
      if (liked.value) return
      
      // 模拟点赞请求
      setTimeout(() => {
        liked.value = true
        ElMessage.success('感谢您的点赞！')
        
        // 将文章ID保存到本地存储
        const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
        likedPosts.push(parseInt(postId.value))
        localStorage.setItem('likedPosts', JSON.stringify(likedPosts))
      }, 500)
    }
    
    // 分享处理
    const handleShare = () => {
      if (navigator.share) {
        navigator.share({
          title: post.value.title,
          text: post.value.summary || '分享一篇好文章',
          url: window.location.href
        }).then(() => {
          ElMessage.success('分享成功！')
        }).catch((error) => {
          console.error('分享失败', error)
        })
      } else {
        // 复制链接到剪贴板
        const url = window.location.href
        navigator.clipboard.writeText(url).then(() => {
          ElMessage.success('链接已复制到剪贴板')
        }).catch(() => {
          ElMessageBox.prompt('复制以下链接分享', '分享', {
            confirmButtonText: '复制',
            cancelButtonText: '取消',
            inputValue: url,
            inputPattern: /.+/,
            inputErrorMessage: '无效链接'
          }).then(({ value }) => {
            ElMessage.success('链接已复制')
          }).catch(() => {})
        })
      }
    }
    
    // 页面加载时获取数据
    onMounted(() => {
      fetchPost()
    })
    
    return {
      post,
      loading,
      liked,
      canEdit,
      renderedContent,
      formatDate,
      handleLike,
      handleShare,
      fetchPost
    }
  }
}
</script>

<style scoped>
.post-detail {
  max-width: 900px;
  margin: 0 auto;
}

.article-cover {
  width: 100%;
  max-height: 500px;
  object-fit: cover;
  border-radius: 8px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.article-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 20px;
  line-height: 1.3;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.article-meta-item {
  margin-right: 20px;
  display: flex;
  align-items: center;
  color: #666;
  font-size: 0.9rem;
}

.article-meta-item i {
  margin-right: 6px;
}

.article-tags {
  margin-bottom: 30px;
}

.tag {
  margin-right: 8px;
}

.article-content {
  line-height: 1.8;
  color: #333;
  font-size: 1.1rem;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  margin-top: 1.5em;
  margin-bottom: 0.8em;
  font-weight: 600;
}

.article-content :deep(p) {
  margin-bottom: 1.2em;
}

.article-content :deep(img) {
  max-width: 100%;
  border-radius: 4px;
  margin: 10px 0;
}

.article-content :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding-left: 16px;
  margin-left: 0;
  color: #666;
}

.article-content :deep(code) {
  background-color: #f5f7fa;
  padding: 3px 5px;
  border-radius: 3px;
  font-family: Consolas, Monaco, 'Courier New', monospace;
}

.article-content :deep(pre) {
  margin: 16px 0;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
  overflow-x: auto;
}

.article-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.article-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.article-content :deep(table th),
.article-content :deep(table td) {
  border: 1px solid #ddd;
  padding: 8px;
}

.article-content :deep(table th) {
  background-color: #f5f7fa;
  font-weight: 600;
}

.article-footer {
  margin-top: 40px;
}

.article-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .article-title {
    font-size: 1.8rem;
  }
  
  .article-content {
    font-size: 1rem;
  }
  
  .article-actions {
    flex-wrap: wrap;
  }
}
</style> 
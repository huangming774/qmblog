<template>
  <div class="user-favorites">
    <div class="page-header">
      <h2 class="page-title">我的收藏</h2>
      
      <div class="filter-bar">
        <el-select v-model="categoryId" placeholder="分类" clearable @change="handleFilterChange">
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
        
        <el-select v-model="tagId" placeholder="标签" clearable @change="handleFilterChange">
          <el-option
            v-for="tag in tags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.id"
          />
        </el-select>
        
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索收藏文章"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="handleSearch">
                <Search />
              </el-icon>
            </template>
          </el-input>
        </div>
      </div>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="8" animated />
    </div>
    
    <div v-else-if="favorites.length === 0" class="empty-container">
      <el-empty description="暂无收藏" :image-size="200">
        <template #description>
          <p>您还没有收藏任何文章</p>
        </template>
        <el-button type="primary" @click="$router.push('/')">
          浏览文章
        </el-button>
      </el-empty>
    </div>
    
    <div v-else class="favorites-grid">
      <div 
        v-for="favorite in favorites" 
        :key="favorite.id" 
        class="favorite-card"
        @click="viewPost(favorite.post)"
      >
        <div class="post-cover" v-if="favorite.post.coverImage">
          <img :src="favorite.post.coverImage" :alt="favorite.post.title" />
        </div>
        <div class="card-content">
          <h3 class="post-title">{{ favorite.post.title }}</h3>
          
          <div class="post-info">
            <div class="post-meta">
              <span class="post-author">
                <el-icon><User /></el-icon>
                {{ favorite.post.author.name }}
              </span>
              <span class="post-date">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(favorite.post.createdAt) }}
              </span>
              <span class="view-count">
                <el-icon><View /></el-icon>
                {{ favorite.post.viewCount }}
              </span>
            </div>
            
            <div class="post-tags" v-if="favorite.post.tags && favorite.post.tags.length > 0">
              <el-tag
                v-for="tag in favorite.post.tags.slice(0, 3)"
                :key="tag.id"
                size="small"
                effect="plain"
                class="post-tag"
              >
                {{ tag.name }}
              </el-tag>
              <span v-if="favorite.post.tags.length > 3" class="more-tags">+{{ favorite.post.tags.length - 3 }}</span>
            </div>
          </div>
          
          <div class="post-summary">{{ favorite.post.summary || truncateContent(favorite.post.content) }}</div>
          
          <div class="favorite-footer">
            <span class="favorite-date">收藏于 {{ formatDate(favorite.createdAt) }}</span>
            
            <div class="favorite-actions">
              <el-popconfirm
                title="确定要取消收藏吗?"
                width="200"
                @confirm.stop="removeFavorite(favorite.id)"
              >
                <template #reference>
                  <el-button 
                    type="danger" 
                    link 
                    size="small" 
                    @click.stop
                  >
                    取消收藏
                  </el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
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
</template>

<script>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, User, Calendar, View } from '@element-plus/icons-vue'

export default {
  name: 'UserFavorites',
  components: {
    Search,
    User,
    Calendar,
    View
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const loading = ref(false)
    const favorites = ref([])
    const categories = ref([])
    const tags = ref([])
    const currentPage = ref(1)
    const pageSize = ref(12)
    const total = ref(0)
    const categoryId = ref('')
    const tagId = ref('')
    const searchKeyword = ref('')
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`
    }
    
    // 截断内容
    const truncateContent = (content) => {
      if (!content) return ''
      // 去除HTML标签
      const plainText = content.replace(/<[^>]+>/g, '')
      if (plainText.length > 100) {
        return plainText.substring(0, 100) + '...'
      }
      return plainText
    }
    
    // 获取分类列表
    const fetchCategories = async () => {
      try {
        const response = await store._vm.$axios.get('/categories')
        categories.value = response.data.data || []
      } catch (error) {
        console.error('获取分类列表失败', error)
      }
    }
    
    // 获取标签列表
    const fetchTags = async () => {
      try {
        const response = await store._vm.$axios.get('/tags')
        tags.value = response.data.data || []
      } catch (error) {
        console.error('获取标签列表失败', error)
      }
    }
    
    // 获取收藏列表
    const fetchFavorites = async () => {
      loading.value = true
      try {
        const response = await store._vm.$axios.get('/user/favorites', {
          params: {
            page: currentPage.value,
            pageSize: pageSize.value,
            categoryId: categoryId.value,
            tagId: tagId.value,
            keyword: searchKeyword.value
          }
        })
        
        favorites.value = response.data.data || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('获取收藏列表失败', error)
        ElMessage.error('获取收藏列表失败，请稍后再试')
      } finally {
        loading.value = false
      }
    }
    
    // 查看文章
    const viewPost = (post) => {
      router.push(`/posts/${post.id}`)
    }
    
    // 取消收藏
    const removeFavorite = async (favoriteId) => {
      try {
        const response = await store._vm.$axios.delete(`/favorites/${favoriteId}`)
        
        if (response.data.success) {
          ElMessage.success('已取消收藏')
          // 更新列表
          favorites.value = favorites.value.filter(favorite => favorite.id !== favoriteId)
          // 如果当前页为空且不是第一页，则返回上一页
          if (favorites.value.length === 0 && currentPage.value > 1) {
            currentPage.value -= 1
            fetchFavorites()
          }
        } else {
          ElMessage.error(response.data.message || '操作失败')
        }
      } catch (error) {
        console.error('取消收藏失败', error)
        ElMessage.error('取消收藏失败，请稍后再试')
      }
    }
    
    // 筛选改变
    const handleFilterChange = () => {
      currentPage.value = 1
      fetchFavorites()
    }
    
    // 搜索
    const handleSearch = () => {
      currentPage.value = 1
      fetchFavorites()
    }
    
    // 切换页码
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchFavorites()
    }
    
    onMounted(() => {
      fetchCategories()
      fetchTags()
      fetchFavorites()
    })
    
    return {
      loading,
      favorites,
      categories,
      tags,
      currentPage,
      pageSize,
      total,
      categoryId,
      tagId,
      searchKeyword,
      formatDate,
      truncateContent,
      viewPost,
      removeFavorite,
      handleFilterChange,
      handleSearch,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.user-favorites {
  min-height: 80vh;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
}

.filter-bar {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  width: 250px;
  margin-left: auto;
}

.search-icon {
  cursor: pointer;
}

.loading-container,
.empty-container {
  padding: 40px 0;
}

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.favorite-card {
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
  background-color: var(--el-bg-color);
  cursor: pointer;
  height: 100%;
}

.favorite-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.post-cover {
  height: 180px;
  overflow: hidden;
}

.post-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.favorite-card:hover .post-cover img {
  transform: scale(1.05);
}

.card-content {
  padding: 16px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.post-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.post-meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.post-author,
.post-date,
.view-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.post-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-top: 5px;
}

.post-tag {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.more-tags {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.post-summary {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-bottom: 16px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.6;
  flex-grow: 1;
}

.favorite-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 12px;
  margin-top: auto;
}

.favorite-date {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .search-box {
    width: 100%;
    margin-left: 0;
  }
  
  .favorites-grid {
    grid-template-columns: 1fr;
  }
}
</style> 
<template>
  <div class="user-comments">
    <div class="page-header">
      <h2 class="page-title">我的评论</h2>
    </div>
    
    <div class="filter-bar">
      <el-select v-model="postId" placeholder="筛选文章" clearable @change="handlePostChange">
        <el-option
          v-for="post in userPosts"
          :key="post.id"
          :label="post.title"
          :value="post.id"
        />
      </el-select>
      
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索评论内容"
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
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="8" animated />
    </div>
    
    <div v-else-if="comments.length === 0" class="empty-container">
      <el-empty description="暂无评论" :image-size="200">
        <template #description>
          <p>您还没有发表任何评论</p>
        </template>
      </el-empty>
    </div>
    
    <div v-else class="comments-list">
      <div 
        v-for="comment in comments" 
        :key="comment.id" 
        class="comment-card"
      >
        <div class="comment-header">
          <div class="comment-meta">
            <router-link :to="`/posts/${comment.postId}`" class="post-link">
              {{ comment.postTitle }}
            </router-link>
            <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
          </div>
          <div class="comment-actions">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="editComment(comment)"
            >
              编辑
            </el-button>
            
            <el-popconfirm
              title="确定要删除此评论吗?"
              width="200"
              @confirm="deleteComment(comment.id)"
            >
              <template #reference>
                <el-button 
                  type="danger" 
                  link 
                  size="small"
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
        
        <div class="comment-content">
          {{ comment.content }}
        </div>
        
        <div v-if="comment.parentId" class="reply-info">
          <el-tag size="small" type="info" effect="plain">回复</el-tag>
          <div class="reply-content">{{ comment.replyTo }}</div>
        </div>
        
        <div v-if="comment.replies && comment.replies.length > 0" class="replies-container">
          <div class="replies-header">
            <span class="replies-count">{{ comment.replies.length }}条回复</span>
          </div>
          <div class="replies-list">
            <div 
              v-for="reply in comment.replies" 
              :key="reply.id" 
              class="reply-item"
            >
              <span class="reply-author">{{ reply.author.name }}:</span>
              <span class="reply-text">{{ reply.content }}</span>
              <span class="reply-date">{{ formatDate(reply.createdAt) }}</span>
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
    
    <!-- 编辑评论对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑评论"
      width="500px"
    >
      <el-form :model="editForm" ref="editFormRef">
        <el-form-item prop="content" :rules="[{ required: true, message: '评论内容不能为空', trigger: 'blur' }]">
          <el-input
            v-model="editForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入评论内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEdit">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

export default {
  name: 'UserComments',
  components: {
    Search
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const loading = ref(false)
    const comments = ref([])
    const userPosts = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    const postId = ref('')
    const searchKeyword = ref('')
    
    const editDialogVisible = ref(false)
    const editForm = ref({
      id: '',
      content: ''
    })
    const editFormRef = ref(null)
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
    }
    
    // 获取用户文章列表
    const fetchUserPosts = async () => {
      try {
        const response = await store._vm.$axios.get('/user/posts', {
          params: {
            page: 1,
            pageSize: 100 // 获取足够多的文章用于筛选
          }
        })
        
        userPosts.value = response.data.data || []
      } catch (error) {
        console.error('获取用户文章列表失败', error)
      }
    }
    
    // 获取评论列表
    const fetchComments = async () => {
      loading.value = true
      try {
        const response = await store._vm.$axios.get('/user/comments', {
          params: {
            page: currentPage.value,
            pageSize: pageSize.value,
            postId: postId.value,
            keyword: searchKeyword.value
          }
        })
        
        comments.value = response.data.data || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('获取评论列表失败', error)
        ElMessage.error('获取评论列表失败，请稍后再试')
      } finally {
        loading.value = false
      }
    }
    
    // 编辑评论
    const editComment = (comment) => {
      editForm.value = {
        id: comment.id,
        content: comment.content
      }
      editDialogVisible.value = true
    }
    
    // 提交编辑
    const submitEdit = async () => {
      if (!editFormRef.value) return
      
      await editFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        try {
          const response = await store._vm.$axios.put(`/comments/${editForm.value.id}`, {
            content: editForm.value.content
          })
          
          if (response.data.success) {
            ElMessage.success('评论已更新')
            editDialogVisible.value = false
            
            // 更新列表
            const index = comments.value.findIndex(c => c.id === editForm.value.id)
            if (index !== -1) {
              comments.value[index].content = editForm.value.content
            }
          } else {
            ElMessage.error(response.data.message || '更新失败')
          }
        } catch (error) {
          console.error('更新评论失败', error)
          ElMessage.error('更新评论失败，请稍后再试')
        }
      })
    }
    
    // 删除评论
    const deleteComment = async (commentId) => {
      try {
        const response = await store._vm.$axios.delete(`/comments/${commentId}`)
        
        if (response.data.success) {
          ElMessage.success('评论已删除')
          // 更新列表
          comments.value = comments.value.filter(comment => comment.id !== commentId)
          // 如果当前页为空且不是第一页，则返回上一页
          if (comments.value.length === 0 && currentPage.value > 1) {
            currentPage.value -= 1
            fetchComments()
          }
        } else {
          ElMessage.error(response.data.message || '删除失败')
        }
      } catch (error) {
        console.error('删除评论失败', error)
        ElMessage.error('删除评论失败，请稍后再试')
      }
    }
    
    // 筛选文章
    const handlePostChange = () => {
      currentPage.value = 1
      fetchComments()
    }
    
    // 搜索
    const handleSearch = () => {
      currentPage.value = 1
      fetchComments()
    }
    
    // 切换页码
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchComments()
    }
    
    onMounted(() => {
      fetchUserPosts()
      fetchComments()
    })
    
    return {
      loading,
      comments,
      userPosts,
      currentPage,
      pageSize,
      total,
      postId,
      searchKeyword,
      editDialogVisible,
      editForm,
      editFormRef,
      formatDate,
      editComment,
      deleteComment,
      handlePostChange,
      handleSearch,
      handlePageChange,
      submitEdit
    }
  }
}
</script>

<style scoped>
.user-comments {
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

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.search-box {
  width: 250px;
}

.search-icon {
  cursor: pointer;
}

.loading-container,
.empty-container {
  padding: 40px 0;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comment-card {
  background-color: var(--el-bg-color);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.comment-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.post-link {
  font-weight: 500;
  color: var(--el-color-primary);
  text-decoration: none;
}

.post-link:hover {
  text-decoration: underline;
}

.comment-date {
  font-size: 12px;
  color: #909399;
}

.comment-actions {
  display: flex;
  gap: 5px;
}

.comment-content {
  line-height: 1.6;
  margin-bottom: 10px;
  word-break: break-word;
}

.reply-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  background-color: var(--el-fill-color-lighter);
  padding: 8px 12px;
  border-radius: 4px;
}

.reply-content {
  font-size: 13px;
  color: #606266;
}

.replies-container {
  margin-top: 10px;
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 10px;
}

.replies-header {
  margin-bottom: 8px;
}

.replies-count {
  font-size: 13px;
  color: #909399;
}

.replies-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.reply-item {
  background-color: var(--el-fill-color-lighter);
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 13px;
}

.reply-author {
  font-weight: 500;
  margin-right: 5px;
}

.reply-text {
  color: #606266;
}

.reply-date {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
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
    gap: 10px;
  }
  
  .search-box {
    width: 100%;
  }
  
  .comment-header {
    flex-direction: column;
    gap: 10px;
  }
  
  .comment-actions {
    align-self: flex-end;
  }
}
</style> 
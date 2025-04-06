<template>
  <div class="comment-section">
    <h3 class="comment-title">评论 ({{ comments.length }})</h3>
    
    <!-- 发表评论表单 -->
    <div class="comment-form" v-if="isAuthenticated">
      <el-form :model="commentForm" :rules="rules" ref="commentFormRef">
        <el-form-item prop="content">
          <el-input
            v-model="commentForm.content"
            type="textarea"
            :rows="3"
            placeholder="写下你的评论..."
            resize="none"
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitComment"
            :loading="loading"
          >
            发表评论
          </el-button>
          <el-button @click="resetForm" v-if="commentForm.content">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="login-tip" v-else>
      <el-alert
        title="请先登录再发表评论"
        type="info"
        show-icon
        :closable="false"
      >
        <template #default>
          <div class="login-buttons">
            <el-button size="small" @click="$router.push('/login')">登录</el-button>
            <el-button size="small" @click="$router.push('/register')">注册</el-button>
          </div>
        </template>
      </el-alert>
    </div>
    
    <!-- 评论列表 -->
    <div class="comment-list" v-if="comments.length > 0">
      <div class="comment-item" v-for="comment in comments" :key="comment.id">
        <div class="comment-header">
          <div class="comment-user">
            <el-avatar :size="40" :src="comment.user?.avatar || defaultAvatar">
              {{ comment.user?.username ? comment.user.username.charAt(0).toUpperCase() : 'U' }}
            </el-avatar>
            <div class="comment-info">
              <div class="comment-username">{{ comment.user?.username || '未知用户' }}</div>
              <div class="comment-time">{{ formatDate(comment.createdAt) }}</div>
            </div>
          </div>
          
          <div class="comment-actions" v-if="canDeleteComment(comment)">
            <el-dropdown @command="handleCommentAction">
              <el-button type="text">
                <i class="fas fa-ellipsis-v"></i>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ type: 'delete', id: comment.id }">
                    <i class="fas fa-trash-alt"></i> 删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
        
        <div class="comment-content">
          {{ comment.content }}
        </div>
        
        <div class="comment-footer">
          <el-button 
            text 
            size="small" 
            @click="replyToComment(comment)"
            v-if="isAuthenticated"
          >
            <i class="fas fa-reply"></i> 回复
          </el-button>
        </div>
        
        <!-- 回复列表 -->
        <div class="replies" v-if="comment.replies && comment.replies.length > 0">
          <div class="reply-item" v-for="reply in comment.replies" :key="reply.id">
            <div class="reply-header">
              <div class="reply-user">
                <el-avatar :size="32" :src="reply.user?.avatar || defaultAvatar">
                  {{ reply.user?.username ? reply.user.username.charAt(0).toUpperCase() : 'U' }}
                </el-avatar>
                <div class="reply-info">
                  <div class="reply-username">{{ reply.user?.username || '未知用户' }}</div>
                  <div class="reply-time">{{ formatDate(reply.createdAt) }}</div>
                </div>
              </div>
              
              <div class="reply-actions" v-if="canDeleteComment(reply)">
                <el-button 
                  type="text" 
                  size="small" 
                  @click="deleteComment(reply.id)"
                >
                  <i class="fas fa-trash-alt"></i>
                </el-button>
              </div>
            </div>
            
            <div class="reply-content">
              {{ reply.content }}
            </div>
          </div>
        </div>
        
        <!-- 回复表单 -->
        <div class="reply-form" v-if="showReplyForm && currentReplyId === comment.id">
          <el-form :model="replyForm" :rules="rules" ref="replyFormRef">
            <el-form-item prop="content">
              <el-input
                v-model="replyForm.content"
                type="textarea"
                :rows="2"
                placeholder="写下你的回复..."
                resize="none"
              />
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary" 
                size="small"
                @click="submitReply"
                :loading="loading"
              >
                提交回复
              </el-button>
              <el-button size="small" @click="cancelReply">取消</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    
    <div class="no-comments" v-else>
      <el-empty description="暂无评论，发表第一条评论吧！"></el-empty>
    </div>
  </div>
</template>

<script>
import { ref, computed, reactive } from 'vue'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'CommentSection',
  props: {
    postId: {
      type: [Number, String],
      required: true
    },
    comments: {
      type: Array,
      default: () => []
    }
  },
  emits: ['refresh-comments'],
  
  setup(props, { emit }) {
    const store = useStore()
    const commentFormRef = ref(null)
    const replyFormRef = ref(null)
    
    // 状态
    const loading = ref(false)
    const showReplyForm = ref(false)
    const currentReplyId = ref(null)
    const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    
    // 表单数据
    const commentForm = reactive({
      content: ''
    })
    
    const replyForm = reactive({
      content: '',
      parentId: null
    })
    
    // 表单验证规则
    const rules = {
      content: [
        { required: true, message: '请输入评论内容', trigger: 'blur' },
        { min: 2, max: 500, message: '评论长度在2到500个字符之间', trigger: 'blur' }
      ]
    }
    
    // 计算属性
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    const currentUser = computed(() => store.state.user)
    
    // 检查是否可以删除评论
    const canDeleteComment = (comment) => {
      if (!isAuthenticated.value || !currentUser.value) return false
      
      return currentUser.value.id === comment.userId || currentUser.value.role === 'admin'
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      return dateString ? new Date(dateString).toLocaleString('zh-CN') : ''
    }
    
    // 重置评论表单
    const resetForm = () => {
      if (commentFormRef.value) {
        commentFormRef.value.resetFields()
      }
    }
    
    // 提交评论
    const submitComment = async () => {
      if (!commentFormRef.value) return
      
      await commentFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        loading.value = true
        
        try {
          const response = await fetch(`http://localhost:8080/api/v1/posts/${props.postId}/comments`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${store.state.token}`
            },
            body: JSON.stringify({
              content: commentForm.content
            })
          })
          
          if (!response.ok) {
            throw new Error('提交评论失败')
          }
          
          ElMessage({
            type: 'success',
            message: '评论发表成功'
          })
          
          resetForm()
          emit('refresh-comments')
        } catch (error) {
          console.error('提交评论出错:', error)
          ElMessage({
            type: 'error',
            message: '提交评论失败，请稍后再试'
          })
        } finally {
          loading.value = false
        }
      })
    }
    
    // 回复评论
    const replyToComment = (comment) => {
      showReplyForm.value = true
      currentReplyId.value = comment.id
      replyForm.parentId = comment.id
      replyForm.content = ''
    }
    
    // 取消回复
    const cancelReply = () => {
      showReplyForm.value = false
      currentReplyId.value = null
      replyForm.parentId = null
      replyForm.content = ''
    }
    
    // 提交回复
    const submitReply = async () => {
      if (!replyFormRef.value) return
      
      await replyFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        loading.value = true
        
        try {
          const response = await fetch(`http://localhost:8080/api/v1/posts/${props.postId}/comments`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${store.state.token}`
            },
            body: JSON.stringify({
              content: replyForm.content,
              parentId: replyForm.parentId
            })
          })
          
          if (!response.ok) {
            throw new Error('提交回复失败')
          }
          
          ElMessage({
            type: 'success',
            message: '回复发表成功'
          })
          
          cancelReply()
          emit('refresh-comments')
        } catch (error) {
          console.error('提交回复出错:', error)
          ElMessage({
            type: 'error',
            message: '提交回复失败，请稍后再试'
          })
        } finally {
          loading.value = false
        }
      })
    }
    
    // 处理评论操作
    const handleCommentAction = ({ type, id }) => {
      if (type === 'delete') {
        deleteComment(id)
      }
    }
    
    // 删除评论
    const deleteComment = async (commentId) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这条评论吗？',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
        )
        
        loading.value = true
        
        const response = await fetch(`http://localhost:8080/api/v1/comments/${commentId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${store.state.token}`
          }
        })
        
        if (!response.ok) {
          throw new Error('删除评论失败')
        }
        
        ElMessage({
          type: 'success',
          message: '评论已删除'
        })
        
        emit('refresh-comments')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除评论出错:', error)
          ElMessage({
            type: 'error',
            message: '删除评论失败，请稍后再试'
          })
        }
      } finally {
        loading.value = false
      }
    }
    
    return {
      commentFormRef,
      replyFormRef,
      loading,
      commentForm,
      replyForm,
      rules,
      isAuthenticated,
      currentUser,
      showReplyForm,
      currentReplyId,
      defaultAvatar,
      canDeleteComment,
      formatDate,
      resetForm,
      submitComment,
      replyToComment,
      cancelReply,
      submitReply,
      handleCommentAction,
      deleteComment
    }
  }
}
</script>

<style scoped>
.comment-section {
  margin-top: 40px;
}

.comment-title {
  font-size: 20px;
  margin-bottom: 20px;
  color: #303133;
  font-weight: 600;
}

.comment-form {
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.login-tip {
  margin-bottom: 30px;
}

.login-buttons {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}

.comment-list {
  margin-top: 20px;
}

.comment-item {
  border-bottom: 1px solid #ebeef5;
  padding: 20px 0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.comment-user {
  display: flex;
  align-items: center;
}

.comment-info {
  margin-left: 10px;
}

.comment-username {
  font-weight: 500;
  font-size: 16px;
  color: #303133;
}

.comment-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.comment-content {
  font-size: 15px;
  line-height: 1.6;
  color: #606266;
  margin: 10px 0;
  white-space: pre-line;
}

.comment-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.replies {
  margin-left: 50px;
  margin-top: 15px;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 10px;
}

.reply-item {
  padding: 10px 0;
  border-bottom: 1px dashed #ebeef5;
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.reply-user {
  display: flex;
  align-items: center;
}

.reply-info {
  margin-left: 10px;
}

.reply-username {
  font-weight: 500;
  font-size: 14px;
  color: #303133;
}

.reply-time {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.reply-content {
  font-size: 14px;
  line-height: 1.6;
  color: #606266;
  margin: 8px 0;
  padding-left: 42px;
  white-space: pre-line;
}

.reply-form {
  margin-top: 15px;
  margin-left: 50px;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 15px;
}

.no-comments {
  margin-top: 30px;
}

@media (max-width: 768px) {
  .comment-form {
    padding: 15px;
  }
  
  .replies {
    margin-left: 20px;
  }
  
  .reply-form {
    margin-left: 20px;
  }
  
  .reply-content {
    padding-left: 20px;
  }
}
</style> 
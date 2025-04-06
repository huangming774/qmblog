import { createStore } from 'vuex'

export default createStore({
  state: {
    token: null,
    user: null,
    isAuthenticated: false,
    posts: [],
    totalPosts: 0,
    currentPost: null,
    loading: false,
    favorites: [],
    notifications: [],
    unreadNotificationsCount: 0,
    theme: localStorage.getItem('theme') || 'light',
    archives: [],
    searchResults: [],
    popularTags: []
  },
  getters: {
    isAdmin(state) {
      return state.user && state.user.role === 'admin'
    },
    isAuthenticated(state) {
      return state.isAuthenticated
    },
    currentUser(state) {
      return state.user
    },
    isDarkMode(state) {
      return state.theme === 'dark'
    },
    hasUnreadNotifications(state) {
      return state.unreadNotificationsCount > 0
    },
    isFavorited: (state) => (postId) => {
      return state.favorites.some(favorite => favorite.postId === postId)
    }
  },
  mutations: {
    // 认证相关
    setUser(state, user) {
      state.user = user
      state.isAuthenticated = true
    },
    setToken(state, token) {
      state.token = token
      // 存储到本地
      localStorage.setItem('token', token)
    },
    clearUserSession(state) {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      state.favorites = []
      state.notifications = []
      state.unreadNotificationsCount = 0
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    // 文章相关
    setPosts(state, { posts, total }) {
      state.posts = posts
      state.totalPosts = total
    },
    setCurrentPost(state, post) {
      state.currentPost = post
    },
    // 加载状态
    setLoading(state, status) {
      state.loading = status
    },
    // 新增mutations
    setFavorites(state, favorites) {
      state.favorites = favorites
    },
    addFavorite(state, favorite) {
      state.favorites.push(favorite)
    },
    removeFavorite(state, postId) {
      state.favorites = state.favorites.filter(fav => fav.postId !== postId)
    },
    setNotifications(state, notifications) {
      state.notifications = notifications
      state.unreadNotificationsCount = notifications.filter(n => !n.read).length
    },
    markNotificationRead(state, notificationId) {
      const notification = state.notifications.find(n => n.id === notificationId)
      if (notification && !notification.read) {
        notification.read = true
        state.unreadNotificationsCount--
      }
    },
    markAllNotificationsRead(state) {
      state.notifications.forEach(notification => {
        notification.read = true
      })
      state.unreadNotificationsCount = 0
    },
    setTheme(state, theme) {
      state.theme = theme
      localStorage.setItem('theme', theme)
      document.documentElement.classList.toggle('dark', theme === 'dark')
    },
    setArchives(state, archives) {
      state.archives = archives
    },
    setSearchResults(state, results) {
      state.searchResults = results
    },
    setPopularTags(state, tags) {
      state.popularTags = tags
    }
  },
  actions: {
    // 用户登录
    async login({ commit }, credentials) {
      try {
        const response = await this._vm.$axios.post('/auth/login', credentials)
        const { token, user } = response.data
        
        commit('setToken', token)
        commit('setUser', user)
        
        // 存储用户信息到本地
        localStorage.setItem('user', JSON.stringify(user))
        
        return { success: true }
      } catch (error) {
        console.error('登录失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '登录失败，请稍后再试' 
        }
      }
    },
    
    // 用户注册
    async register({ commit }, userData) {
      try {
        const response = await this._vm.$axios.post('/auth/register', userData)
        const { token, user } = response.data
        
        commit('setToken', token)
        commit('setUser', user)
        
        // 存储用户信息到本地
        localStorage.setItem('user', JSON.stringify(user))
        
        return { success: true }
      } catch (error) {
        console.error('注册失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '注册失败，请稍后再试' 
        }
      }
    },
    
    // 用户登出
    logout({ commit }) {
      commit('clearUserSession')
    },
    
    // 获取文章列表
    async fetchPosts({ commit }, { page = 1, pageSize = 10, status = 'published', tag = '' }) {
      commit('setLoading', true)
      try {
        const params = { page, pageSize, status }
        if (tag) params.tag = tag
        
        const response = await this._vm.$axios.get('/posts', { params })
        
        commit('setPosts', {
          posts: response.data.data,
          total: response.data.total
        })
        
        return { success: true }
      } catch (error) {
        console.error('获取文章列表失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取文章列表失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },
    
    // 获取文章详情
    async fetchPost({ commit }, postId) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.get(`/posts/${postId}`)
        commit('setCurrentPost', response.data)
        return { success: true }
      } catch (error) {
        console.error('获取文章详情失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取文章详情失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },
    
    // 创建文章
    async createPost({ commit }, postData) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.post('/posts', postData)
        return { 
          success: true,
          postId: response.data.id
        }
      } catch (error) {
        console.error('创建文章失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '创建文章失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },
    
    // 更新文章
    async updatePost({ commit }, { postId, postData }) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.put(`/posts/${postId}`, postData)
        return { 
          success: true,
          post: response.data
        }
      } catch (error) {
        console.error('更新文章失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '更新文章失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },
    
    // 删除文章
    async deletePost({ commit }, postId) {
      commit('setLoading', true)
      try {
        await this._vm.$axios.delete(`/posts/${postId}`)
        return { success: true }
      } catch (error) {
        console.error('删除文章失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '删除文章失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 新增actions
    // 更新用户资料
    async updateUserProfile({ commit }, userData) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.put('/user/profile', userData)
        commit('setUser', response.data)
        localStorage.setItem('user', JSON.stringify(response.data))
        return { success: true }
      } catch (error) {
        console.error('更新用户资料失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '更新用户资料失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 修改密码
    async changePassword({ commit }, { oldPassword, newPassword }) {
      commit('setLoading', true)
      try {
        await this._vm.$axios.put('/user/password', { oldPassword, newPassword })
        return { success: true }
      } catch (error) {
        console.error('修改密码失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '修改密码失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 获取收藏列表
    async fetchFavorites({ commit }) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.get('/user/favorites')
        commit('setFavorites', response.data)
        return { success: true }
      } catch (error) {
        console.error('获取收藏列表失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取收藏列表失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 添加收藏
    async addFavorite({ commit }, postId) {
      try {
        const response = await this._vm.$axios.post('/user/favorites', { postId })
        commit('addFavorite', response.data)
        return { success: true }
      } catch (error) {
        console.error('添加收藏失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '添加收藏失败' 
        }
      }
    },

    // 取消收藏
    async removeFavorite({ commit }, favoriteId) {
      try {
        await this._vm.$axios.delete(`/user/favorites/${favoriteId}`)
        commit('removeFavorite', favoriteId)
        return { success: true }
      } catch (error) {
        console.error('取消收藏失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '取消收藏失败' 
        }
      }
    },

    // 获取通知列表
    async fetchNotifications({ commit }) {
      try {
        const response = await this._vm.$axios.get('/user/notifications')
        commit('setNotifications', response.data)
        return { success: true }
      } catch (error) {
        console.error('获取通知列表失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取通知列表失败' 
        }
      }
    },

    // 标记通知为已读
    async markNotificationRead({ commit }, notificationId) {
      try {
        await this._vm.$axios.put(`/user/notifications/${notificationId}/read`)
        commit('markNotificationRead', notificationId)
        return { success: true }
      } catch (error) {
        console.error('标记通知失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '标记通知失败' 
        }
      }
    },

    // 标记所有通知为已读
    async markAllNotificationsRead({ commit }) {
      try {
        await this._vm.$axios.put('/user/notifications/read-all')
        commit('markAllNotificationsRead')
        return { success: true }
      } catch (error) {
        console.error('标记所有通知失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '标记所有通知失败' 
        }
      }
    },

    // 切换主题
    toggleTheme({ commit, state }) {
      const newTheme = state.theme === 'light' ? 'dark' : 'light'
      commit('setTheme', newTheme)
    },

    // 搜索文章
    async searchPosts({ commit }, { keyword, page = 1, pageSize = 10 }) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.get('/posts/search', { 
          params: { keyword, page, pageSize } 
        })
        commit('setSearchResults', response.data)
        return { 
          success: true, 
          total: response.data.total
        }
      } catch (error) {
        console.error('搜索文章失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '搜索文章失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 获取归档数据
    async fetchArchives({ commit }) {
      commit('setLoading', true)
      try {
        const response = await this._vm.$axios.get('/posts/archives')
        commit('setArchives', response.data)
        return { success: true }
      } catch (error) {
        console.error('获取归档数据失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取归档数据失败' 
        }
      } finally {
        commit('setLoading', false)
      }
    },

    // 获取热门标签
    async fetchPopularTags({ commit }) {
      try {
        const response = await this._vm.$axios.get('/tags/popular')
        commit('setPopularTags', response.data)
        return { success: true }
      } catch (error) {
        console.error('获取热门标签失败', error)
        return { 
          success: false, 
          message: error.response?.data?.error || '获取热门标签失败' 
        }
      }
    }
  }
}) 
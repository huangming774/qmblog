import { createRouter, createWebHistory } from 'vue-router'
import store from '@/store'
import UserLayout from '@/views/user/UserLayout.vue'
import UserProfile from '@/views/user/UserProfile.vue'
import UserPosts from '@/views/user/UserPosts.vue'
import UserComments from '@/views/user/UserComments.vue'
import UserFavorites from '@/views/user/UserFavorites.vue'
import UserSettings from '@/views/user/UserSettings.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/posts/:id',
    name: 'PostDetail',
    component: () => import('../views/posts/PostDetail.vue')
  },
  // 用户中心路由
  {
    path: '/user',
    component: UserLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'profile',
        name: 'UserProfile',
        component: UserProfile,
        meta: { title: '个人资料' }
      },
      {
        path: 'posts',
        name: 'UserPosts',
        component: UserPosts,
        meta: { title: '我的文章' }
      },
      {
        path: 'comments',
        name: 'UserComments',
        component: UserComments,
        meta: { title: '我的评论' }
      },
      {
        path: 'favorites',
        name: 'UserFavorites',
        component: UserFavorites,
        meta: { title: '我的收藏' }
      },
      {
        path: 'settings',
        name: 'UserSettings',
        component: UserSettings,
        meta: { title: '账号设置' }
      },
      {
        path: '',
        redirect: '/user/profile'
      }
    ]
  },
  // 搜索路由
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/search/SearchResults.vue')
  },
  // 归档路由
  {
    path: '/archive',
    name: 'Archive',
    component: () => import('../views/archive/Archive.vue')
  },
  // 标签路由
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('../views/archive/Tags.vue')
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/admin/Dashboard.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'posts',
        name: 'AdminPosts',
        component: () => import('../views/admin/Posts.vue')
      },
      {
        path: 'posts/create',
        name: 'CreatePost',
        component: () => import('../views/admin/EditPost.vue')
      },
      {
        path: 'posts/edit/:id',
        name: 'EditPost',
        component: () => import('../views/admin/EditPost.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  // 滚动行为
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 全局导航守卫
router.beforeEach((to, from, next) => {
  // 需要登录的路由
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } 
  // 只允许未登录用户访问的路由(如登录和注册页)
  else if (to.matched.some(record => record.meta.guest)) {
    if (store.state.isAuthenticated) {
      next({ path: '/' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router 
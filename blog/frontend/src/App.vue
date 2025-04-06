<template>
  <div id="app" class="flex flex-col min-h-screen font-sans">
    <el-config-provider :locale="zhCn">
      <Navbar />
      <div class="main-container flex-1 px-4 py-6 md:py-8 max-w-7xl mx-auto w-full">
        <router-view v-slot="{ Component }">
          <transition 
            name="fade" 
            mode="out-in"
            @before-leave="beforeLeave"
            @enter="enter"
            @after-enter="afterEnter"
          >
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
      <Footer />
    </el-config-provider>
  </div>
</template>

<script>
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import Navbar from '@/components/layout/Navbar.vue'
import Footer from '@/components/layout/Footer.vue'

export default {
  name: 'App',
  components: {
    Navbar,
    Footer
  },
  data() {
    return {
      zhCn
    }
  },
  methods: {
    beforeLeave(el) {
      // 在页面离开前保存滚动位置
      document.body.classList.add('disable-scroll')
    },
    enter(el, done) {
      // 淡入动画
      document.body.style.opacity = 0
      setTimeout(() => {
        document.body.style.opacity = 1
        done()
      }, 50)
    },
    afterEnter() {
      // 恢复滚动并滚动到顶部
      document.body.classList.remove('disable-scroll')
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  },
  mounted() {
    // 检查本地存储中的用户会话
    const token = localStorage.getItem('token')
    const user = localStorage.getItem('user')
    
    if (token && user) {
      try {
        this.$store.commit('setUser', JSON.parse(user))
        this.$store.commit('setToken', token)
      } catch (e) {
        console.error('解析用户信息失败', e)
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      }
    }

    // 初始化页面加载动画
    document.body.classList.add('page-loaded')
  }
}
</script>

<style>
/* 基础样式将由Tailwind CSS处理 */

/* 页面过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 页面加载动画 */
.page-loaded {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* 禁用滚动时的样式 */
.disable-scroll {
  overflow: hidden;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style> 
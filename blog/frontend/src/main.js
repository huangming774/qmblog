import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'
import moment from 'moment'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

// 导入新添加的库
import AOS from 'aos'
import 'aos/dist/aos.css'
import { MotionPlugin } from '@vueuse/motion'
import 'swiper/css'
import 'swiper/css/navigation'
import 'swiper/css/pagination'

// 导入Tailwind CSS主样式
import './assets/main.css'

// 设置中文语言
moment.locale('zh-cn')

// 初始化AOS动画库
AOS.init({
    duration: 800,
    easing: 'ease-in-out',
    once: true
})

// 配置Marked
marked.setOptions({
    renderer: new marked.Renderer(),
    highlight: function(code, lang) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
    },
    langPrefix: 'hljs language-',
    pedantic: false,
    gfm: true,
    breaks: false,
    sanitize: false,
    smartypants: false,
    xhtml: false
});

// 配置Axios
axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8080/api/v1';
axios.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 添加响应拦截器处理401状态
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      store.commit('clearUserSession')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

const app = createApp(App)

// 全局过滤器
app.config.globalProperties.$filters = {
    formatDate(date) {
        return moment(date).format('YYYY-MM-DD HH:mm');
    },
    markdown(content) {
        return marked(content);
    }
}

app.config.globalProperties.$axios = axios
app.config.globalProperties.$moment = moment

// 使用插件
app.use(store)
   .use(router)
   .use(ElementPlus)
   .use(MotionPlugin) // 添加动画插件
   .mount('#app') 
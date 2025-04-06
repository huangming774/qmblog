<template>
  <footer class="footer">
    <div class="footer-content">
      <div class="footer-section">
        <h3>博客系统</h3>
        <p>一个基于Vue和Go的现代化博客系统。</p>
      </div>
      
      <div class="footer-section">
        <h3>链接</h3>
        <ul>
          <li><router-link to="/">首页</router-link></li>
          <li><router-link to="/login">登录</router-link></li>
          <li><router-link to="/register">注册</router-link></li>
        </ul>
      </div>
      
      <div class="footer-section">
        <h3>技术栈</h3>
        <ul>
          <li>前端: Vue 3, Vuex, Router</li>
          <li>后端: Go, Gin, GORM</li>
          <li>数据库: PostgreSQL, Redis</li>
        </ul>
      </div>
      
      <div class="footer-section">
        <h3>联系我们</h3>
        <div class="social-links">
          <a href="#" target="_blank"><i class="fab fa-github"></i></a>
          <a href="#" target="_blank"><i class="fab fa-twitter"></i></a>
          <a href="#" target="_blank"><i class="fab fa-weixin"></i></a>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p>&copy; {{ currentYear }} 博客系统 - 保留所有权利</p>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'AppFooter',
  data() {
    return {
      currentYear: new Date().getFullYear()
    }
  }
}
</script>

<style scoped>
.footer {
  background-color: #2c3e50;
  color: #fff;
  padding: 40px 20px 20px;
  margin-top: 40px;
}

.footer-content {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
}

.footer-section {
  flex: 1;
  min-width: 200px;
  margin-bottom: 20px;
  padding: 0 10px;
}

.footer-section h3 {
  margin-bottom: 15px;
  font-size: 18px;
  color: #fff;
}

.footer-section ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-section ul li {
  margin-bottom: 8px;
}

.footer-section a {
  color: #ddd;
  text-decoration: none;
  transition: color 0.3s;
}

.footer-section a:hover {
  color: #fff;
}

.social-links {
  display: flex;
  gap: 12px;
}

.social-links a {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  transition: background-color 0.3s;
}

.social-links a:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.social-links i {
  font-size: 18px;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 20px;
  margin-top: 20px;
  text-align: center;
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .footer-content {
    flex-direction: column;
  }
  
  .footer-section {
    margin-bottom: 30px;
  }
}
</style> 
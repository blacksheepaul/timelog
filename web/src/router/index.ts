import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 页面组件懒加载
const Home = () => import('@/views/Home.vue')
const TimeLog = () => import('@/views/TimeLog.vue')
const Tags = () => import('@/views/Tags.vue')
const Statistics = () => import('@/views/Statistics.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Dashboard'
    }
  },
  {
    path: '/timelogs',
    name: 'TimeLog',
    component: TimeLog,
    meta: {
      title: 'Time Logs'
    }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: Tags,
    meta: {
      title: 'Tag Management'
    }
  },
  {
    path: '/statistics',
    name: 'Statistics',
    component: Statistics,
    meta: {
      title: 'Statistics'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫 - 设置页面标题
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - TimeLog`
  }
  next()
})

export default router
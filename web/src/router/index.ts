import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 页面组件懒加载
const Home = () => import('@/views/Home.vue')
const TimeLog = () => import('@/views/TimeLog.vue')
const Tasks = () => import('@/views/Tasks.vue')
const Tags = () => import('@/views/Tags.vue')
const Statistics = () => import('@/views/Statistics.vue')
// @ts-ignore
const Constraints = () => import('@/views/Constraints.vue') as Promise<{ default: any }>

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Dashboard',
    },
  },
  {
    path: '/timelogs',
    name: 'TimeLog',
    component: TimeLog,
    meta: {
      title: 'Time Logs',
    },
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: Tasks,
    meta: {
      title: 'Task Management',
    },
  },
  {
    path: '/tags',
    name: 'Tags',
    component: Tags,
    meta: {
      title: 'Tag Management',
    },
  },
  {
    path: '/statistics',
    name: 'Statistics',
    component: Statistics,
    meta: {
      title: 'Statistics',
    },
  },
  {
    path: '/constraints',
    name: 'Constraints',
    component: Constraints,
    meta: {
      title: '约束',
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})

// 路由守卫 - 设置页面标题
router.beforeEach((to, _from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - TimeLog`
  }
  next()
})

export default router

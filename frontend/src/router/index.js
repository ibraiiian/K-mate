import { createRouter, createWebHistory } from 'vue-router'

// Lazy load components
const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardView.vue')
const POSView = () => import('@/views/POSView.vue')
const InventoryView = () => import('@/views/InventoryView.vue')

const TransactionView = () => import('@/views/TransactionsView.vue')
const ReportView = () => import('@/views/ReportView.vue')
const AuditLogView = () => import('@/views/AuditLogView.vue')

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: LoginView,
        meta: { requiresAuth: false, layout: 'blank' },
    },
    {
        path: '/',
        redirect: '/dashboard',
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: DashboardView,
        meta: { requiresAuth: true, roles: ['owner'] },
    },
    {
        path: '/pos',
        name: 'POS',
        component: POSView,
        meta: { requiresAuth: true, roles: ['owner', 'admin', 'cashier'] },
    },
    {
        path: '/inventory',
        name: 'Inventory',
        component: InventoryView,
        meta: { requiresAuth: true, roles: ['owner', 'admin'] },
    },

    {
        path: '/transactions',
        name: 'Transactions',
        component: TransactionView,
        meta: { requiresAuth: true, roles: ['owner', 'admin', 'cashier'] },
    },
    {
        path: '/reports',
        name: 'Reports',
        component: ReportView,
        meta: { requiresAuth: true, roles: ['owner', 'admin'] },
    },
    {
        path: '/audit-logs',
        name: 'AuditLogs',
        component: AuditLogView,
        meta: { requiresAuth: true, roles: ['owner'] },
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/login',
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

// Navigation Guard
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const user = JSON.parse(localStorage.getItem('user'))

    // If route requires auth
    if (to.meta.requiresAuth) {
        if (!token) {
            // No token, redirect to login
            return next({ name: 'Login' })
        }

        // Check role-based access
        if (to.meta.roles && user) {
            if (!to.meta.roles.includes(user.role)) {
                // User doesn't have permission
                alert('Anda tidak memiliki akses ke halaman ini')
                return next({ name: 'POS' }) // Redirect to POS as default
            }
        }
    }

    // If already logged in and trying to access login page
    if (to.name === 'Login' && token) {
        return next({ name: 'Dashboard' })
    }

    next()
})

export default router

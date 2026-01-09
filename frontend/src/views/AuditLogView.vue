<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-slate-800">Audit Log</h1>
        <p class="text-slate-500 text-sm">Pantau aktivitas pengguna sistem</p>
      </div>

      <!-- Filters -->
      <div class="flex flex-col md:flex-row gap-4 mb-6">
        <div class="flex-1 flex gap-3">
          <!-- Action Filter -->
          <select 
            v-model="filterAction" 
            @change="fetchLogs"
            class="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-blue-600 bg-white"
          >
            <option value="">Semua Aksi</option>
            <option v-for="action in actionTypes" :key="action" :value="action">
              {{ formatActionLabel(action) }}
            </option>
          </select>

          <!-- Refresh Button -->
          <BaseButton 
            variant="outline"
            size="icon"
            @click="fetchLogs" 
            :disabled="loading"
            title="Refresh Data"
          >
            <RefreshCcw class="w-4 h-4" :class="{'animate-spin': loading}" />
          </BaseButton>
        </div>

        <!-- Stats -->
        <div class="flex items-center gap-2 text-sm text-slate-500">
        <!-- Stats -->
        <div class="flex items-center gap-2 text-sm text-slate-500">
          <List class="w-4 h-4" />
          <span>Total: <strong class="text-slate-700">{{ totalLogs }}</strong> log</span>
        </div>
          <span>Total: <strong class="text-slate-700">{{ totalLogs }}</strong> log</span>
        </div>
      </div>

      <!-- Data Table -->
      <div class="bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden flex-1 flex flex-col">
        
        <!-- Loading State -->
        <div v-if="loading" class="flex-1 flex items-center justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>

        <!-- Empty State -->
        <div v-else-if="logs.length === 0" class="flex-1 flex flex-col items-center justify-center text-slate-400">
          <ClipboardList class="w-12 h-12 mb-2 opacity-50" />
          <p>Belum ada aktivitas tercatat</p>
        </div>

        <!-- Table -->
        <div v-else class="overflow-x-auto flex-1">
          <table class="w-full data-table">
            <thead>
              <tr>
                <th class="w-40">Waktu</th>
                <th>User</th>
                <th>Aksi</th>
                <th>Target</th>
                <th class="text-center w-20">Detail</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr 
                v-for="log in logs" 
                :key="log.id"
                class="hover:bg-gray-50 cursor-pointer"
                @click="openDetailModal(log)"
              >
                <td class="text-slate-500 text-xs">{{ formatDate(log.timestamp) }}</td>
                <td>
                  <div class="flex items-center gap-2">
                    <div class="w-7 h-7 rounded-full bg-slate-200 flex items-center justify-center text-slate-600 text-xs font-bold">
                      {{ getInitials(log.user_name) }}
                    </div>
                    <span class="font-medium text-slate-700">{{ log.user_name || 'System' }}</span>
                  </div>
                </td>
                <td>
                  <span :class="['action-badge', getActionClass(log.action)]">
                    <component :is="getActionIcon(log.action)" class="w-3 h-3 mr-1" />
                    {{ formatActionLabel(log.action) }}
                  </span>
                </td>
                <td class="text-slate-600">
                  <span class="font-mono text-xs bg-gray-100 px-2 py-1 rounded">{{ log.target_table }}</span>
                  <span v-if="log.target_id" class="text-slate-400 ml-1">#{{ log.target_id }}</span>
                </td>
                <td class="text-center">
                  <button class="text-slate-400 hover:text-blue-600 transition">
                    <Eye class="w-5 h-5" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="logs.length > 0" class="px-6 py-4 border-t border-gray-200 flex justify-between items-center bg-gray-50">
          <p class="text-xs text-slate-500">
            Halaman <span class="font-bold text-slate-700">{{ page }}</span> dari {{ totalPages }}
          </p>
          <div class="flex gap-1">
            <button 
              @click="prevPage" 
              :disabled="page <= 1"
              class="px-3 py-1 border border-gray-300 rounded bg-white text-slate-600 hover:bg-gray-50 text-xs font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Sebelumnya
            </button>
            <button 
              @click="nextPage" 
              :disabled="page >= totalPages"
              class="px-3 py-1 border border-gray-300 rounded bg-white text-slate-600 hover:bg-gray-50 text-xs font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Selanjutnya
            </button>
          </div>
        </div>
      </div>

      <!-- DETAIL MODAL -->
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="closeModal"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-white rounded-lg shadow-xl w-[700px] max-h-[85vh] overflow-hidden border border-gray-200 animate-fadeSlideIn">
          <!-- Header -->
          <div class="bg-slate-50 px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <div>
              <h3 class="font-bold text-slate-800 text-lg flex items-center gap-2">
                <FileText class="w-5 h-5" />
                Detail Aktivitas
              </h3>
              <p class="text-xs text-slate-500 mt-0.5">ID: {{ selectedLog?.id }}</p>
            </div>
            <button @click="closeModal" class="text-slate-400 hover:text-red-600 transition">
              <X class="w-6 h-6" />
            </button>
          </div>

          <!-- Body -->
          <div class="p-6 overflow-y-auto max-h-[60vh]">
            <!-- Info Grid -->
            <div class="grid grid-cols-2 gap-4 mb-6 text-sm">
              <div>
                <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Waktu</p>
                <p class="font-medium text-slate-800">{{ formatDate(selectedLog?.timestamp) }}</p>
              </div>
              <div>
                <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">User</p>
                <p class="font-medium text-slate-800">{{ selectedLog?.user_name || 'System' }}</p>
              </div>
              <div>
                <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Aksi</p>
                <span :class="['action-badge', getActionClass(selectedLog?.action)]">
                  {{ formatActionLabel(selectedLog?.action) }}
                </span>
              </div>
              <div>
                <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Target</p>
                <p class="font-medium text-slate-800">
                  {{ selectedLog?.target_table }} 
                  <span v-if="selectedLog?.target_id" class="text-slate-500">#{{ selectedLog?.target_id }}</span>
                </p>
              </div>
              <div class="col-span-2">
                <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">IP Address</p>
                <p class="font-mono text-slate-600 text-xs">{{ selectedLog?.ip_address || '-' }}</p>
              </div>
            </div>

            <!-- Before / After Comparison -->
            <div v-if="selectedLog?.old_value || selectedLog?.new_value" class="space-y-4">
              <h4 class="font-bold text-slate-700 border-b border-gray-200 pb-2">Perubahan Data</h4>
              
              <div class="grid grid-cols-2 gap-4">
                <!-- Before -->
                <div v-if="selectedLog?.old_value">
                  <p class="text-xs font-bold text-red-600 uppercase mb-2 flex items-center gap-1">
                    <MinusCircle class="w-4 h-4" /> Sebelum
                  </p>
                  <pre class="bg-red-50 border border-red-200 rounded-lg p-3 text-xs font-mono overflow-auto max-h-48 text-red-800">{{ formatJSON(selectedLog?.old_value) }}</pre>
                </div>
                <div v-else class="flex items-center justify-center text-slate-400 text-sm">
                  <span>Tidak ada data sebelumnya</span>
                </div>

                <!-- After -->
                <div v-if="selectedLog?.new_value">
                  <p class="text-xs font-bold text-green-600 uppercase mb-2 flex items-center gap-1">
                    <PlusCircle class="w-4 h-4" /> Sesudah
                  </p>
                  <pre class="bg-green-50 border border-green-200 rounded-lg p-3 text-xs font-mono overflow-auto max-h-48 text-green-800">{{ formatJSON(selectedLog?.new_value) }}</pre>
                </div>
                <div v-else class="flex items-center justify-center text-slate-400 text-sm">
                  <span>Tidak ada data baru</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex justify-end">
            <button 
              @click="closeModal"
              class="px-4 py-2 text-sm font-medium bg-gray-200 hover:bg-gray-300 text-slate-700 rounded-md transition"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>

import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import api from '@/api/axiosInstance'
import BaseButton from '@/components/common/BaseButton.vue'
import { 
  RefreshCcw, 
  List, 
  ClipboardList, 
  Eye, 
  FileText, 
  X, 
  MinusCircle, 
  PlusCircle,
  LogIn,
  AlertTriangle,
  LogOut,
  Plus,
  Pencil,
  Trash2,
  ShoppingCart,
  CornerUpLeft,
  DollarSign,
  Info
} from 'lucide-vue-next'

// State
const loading = ref(false)
const logs = ref([])
const totalLogs = ref(0)
const page = ref(1)
const limit = ref(20)
const filterAction = ref('')
const actionTypes = ref([])

// Modal State
const showDetailModal = ref(false)
const selectedLog = ref(null)

// Computed
const totalPages = computed(() => Math.ceil(totalLogs.value / limit.value) || 1)

// Fetch Logs
const fetchLogs = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: page.value,
      limit: limit.value,
    })
    if (filterAction.value) {
      params.append('action', filterAction.value)
    }

    const response = await api.get(`/audit-logs?${params}`)
    const data = response.data?.data
    logs.value = data?.logs || []
    totalLogs.value = data?.total || 0
  } catch (error) {
    console.error('Error fetching audit logs:', error)
  } finally {
    loading.value = false
  }
}

// Fetch Action Types
const fetchActionTypes = async () => {
  try {
    const response = await api.get('/audit-logs/actions')
    actionTypes.value = response.data?.data || []
  } catch (error) {
    console.error('Error fetching action types:', error)
  }
}

// Pagination
const prevPage = () => {
  if (page.value > 1) {
    page.value--
    fetchLogs()
  }
}

const nextPage = () => {
  if (page.value < totalPages.value) {
    page.value++
    fetchLogs()
  }
}

// Modal Functions
const openDetailModal = (log) => {
  selectedLog.value = log
  showDetailModal.value = true
}

const closeModal = () => {
  showDetailModal.value = false
  selectedLog.value = null
}

// Helpers
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getInitials = (name) => {
  if (!name) return '?'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const formatActionLabel = (action) => {
  const labels = {
    'LOGIN': 'Login',
    'LOGIN_FAILED': 'Login Gagal',
    'LOGOUT': 'Logout',
    'CREATE': 'Tambah',
    'UPDATE': 'Ubah',
    'DELETE': 'Hapus',
    'CHECKOUT': 'Checkout',
    'RETURN': 'Retur',
    'PRICE_CHANGE': 'Ubah Harga'
  }
  return labels[action] || action
}

const getActionClass = (action) => {
  const classes = {
    'LOGIN': 'action-info',
    'LOGIN_FAILED': 'action-delete',
    'LOGOUT': 'action-info',
    'CREATE': 'action-create',
    'UPDATE': 'action-update',
    'DELETE': 'action-delete',
    'CHECKOUT': 'action-create',
    'RETURN': 'action-update',
    'PRICE_CHANGE': 'action-update'
  }
  return classes[action] || 'action-info'
}

const getActionIcon = (action) => {
  const icons = {
    'LOGIN': LogIn,
    'LOGIN_FAILED': AlertTriangle,
    'LOGOUT': LogOut,
    'CREATE': Plus,
    'UPDATE': Pencil,
    'DELETE': Trash2,
    'CHECKOUT': ShoppingCart,
    'RETURN': CornerUpLeft,
    'PRICE_CHANGE': DollarSign
  }
  return icons[action] || Info
}

const formatJSON = (data) => {
  if (!data) return ''
  try {
    return JSON.stringify(data, null, 2)
  } catch (e) {
    return String(data)
  }
}

// Lifecycle
onMounted(() => {
  fetchActionTypes()
  fetchLogs()
})
</script>

<style scoped>
/* Table Styles */
.data-table th {
  background-color: #F8FAFC;
  color: #64748B;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 12px 20px;
  border-bottom: 1px solid #E2E8F0;
  text-align: left;
}
.data-table td {
  padding: 14px 20px;
  font-size: 0.875rem;
}

/* Action Badges */
.action-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}

.action-create {
  background-color: #DCFCE7;
  color: #166534;
  border: 1px solid #BBF7D0;
}

.action-update {
  background-color: #FEF3C7;
  color: #92400E;
  border: 1px solid #FDE68A;
}

.action-delete {
  background-color: #FEE2E2;
  color: #991B1B;
  border: 1px solid #FECACA;
}

.action-info {
  background-color: #E0E7FF;
  color: #3730A3;
  border: 1px solid #C7D2FE;
}

/* Animation */
@keyframes fadeSlideIn { 
  from { opacity: 0; transform: translateY(10px) scale(0.98); } 
  to { opacity: 1; transform: translateY(0) scale(1); } 
}
.animate-fadeSlideIn {
  animation: fadeSlideIn 0.2s ease-out;
}
</style>

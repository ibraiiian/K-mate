<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-slate-800">Laporan Keuangan</h1>
        <p class="text-slate-500 text-sm">Analisa pendapatan dan pengeluaran berdasarkan periode</p>
      </div>

      <!-- 1. DATE FILTER -->
      <div class="bg-white border border-gray-200 rounded-lg p-4 mb-8 shadow-sm">
        <div class="flex flex-col md:flex-row items-center gap-4">
          <div class="flex items-center gap-2 w-full md:w-auto">
            <label class="text-sm font-medium text-slate-600">Dari:</label>
            <input 
              v-model="startDate" 
              type="date" 
              class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-600"
            >
          </div>
          <div class="flex items-center gap-2 w-full md:w-auto">
            <label class="text-sm font-medium text-slate-600">Sampai:</label>
            <input 
              v-model="endDate" 
              type="date" 
              class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-600"
            >
          </div>
          <BaseButton 
            @click="fetchReport" 
            :disabled="loading"
            class="bg-blue-600 hover:bg-blue-700"
          >
            <Search class="w-4 h-4 mr-2" :class="{'animate-spin': loading}" />
            Tampilkan Laporan
          </BaseButton>
        </div>
      </div>

      <!-- 2. SUMMARY CARDS -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <!-- Omzet -->
        <div class="bg-white border border-gray-200 p-5 rounded-lg shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs text-slate-500 font-bold uppercase">Omzet (Revenue)</p>
            <div class="bg-blue-50 p-2 rounded-md text-blue-600"><DollarSign class="w-6 h-6" /></div>
          </div>
          <p class="text-2xl font-bold text-slate-800 font-mono">{{ formatRupiah(report.total_revenue) }}</p>
          <p class="text-xs text-slate-400 mt-1">{{ report.transaction_count }} transaksi</p>
        </div>

        <!-- HPP -->
        <div class="bg-white border border-gray-200 p-5 rounded-lg shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs text-slate-500 font-bold uppercase">HPP (COGS)</p>
            <div class="bg-orange-50 p-2 rounded-md text-orange-600"><Barcode class="w-6 h-6" /></div>
          </div>
          <p class="text-2xl font-bold text-orange-600 font-mono">{{ formatRupiah(report.total_hpp) }}</p>
          <p class="text-xs text-slate-400 mt-1">Harga Pokok Penjualan</p>
        </div>

        <!-- Total Expense -->
        <div class="bg-white border border-gray-200 p-5 rounded-lg shadow-sm">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs text-slate-500 font-bold uppercase">Pengeluaran</p>
            <div class="bg-red-50 p-2 rounded-md text-red-600"><ArrowUpRight class="w-6 h-6" /></div>
          </div>
          <p class="text-2xl font-bold text-red-600 font-mono">{{ formatRupiah(report.total_expense) }}</p>
          <p class="text-xs text-slate-400 mt-1">Beban Operasional</p>
        </div>

        <!-- Net Profit -->
        <div class="bg-slate-800 text-white p-5 rounded-lg shadow-sm border border-slate-900">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs text-slate-400 font-bold uppercase">Laba Bersih</p>
            <div class="bg-slate-700 p-2 rounded-md text-green-400"><TrendingUp class="w-6 h-6" /></div>
          </div>
          <p class="text-2xl font-bold font-mono" :class="parseFloat(report.net_profit) >= 0 ? 'text-green-400' : 'text-red-400'">
            {{ formatRupiah(report.net_profit) }}
          </p>
          <p class="text-xs text-slate-400 mt-1">Revenue - HPP - Expense</p>
        </div>
      </div>

      <!-- 3. EXPENSE BREAKDOWN TABLE -->
      <div class="bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden">
        <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
          <h3 class="font-bold text-slate-800">Rincian Pengeluaran per Kategori</h3>
        </div>
        
        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center p-10">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>

        <!-- Empty State -->
        <div v-else-if="!report.expense_breakdown || report.expense_breakdown.length === 0" class="p-10 text-center text-slate-400">
          <ClipboardList class="w-12 h-12 mb-2 mx-auto opacity-50" />
          <p>Tidak ada data pengeluaran untuk periode ini</p>
        </div>

        <!-- Table -->
        <table v-else class="w-full">
          <thead class="bg-gray-50 text-left">
            <tr>
              <th class="px-4 py-3 text-xs font-bold text-slate-500 uppercase">Kategori</th>
              <th class="px-4 py-3 text-xs font-bold text-slate-500 uppercase text-center">Jumlah Transaksi</th>
              <th class="px-4 py-3 text-xs font-bold text-slate-500 uppercase text-right">Total</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="(item, index) in report.expense_breakdown" :key="index" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-slate-700 capitalize">{{ item.category }}</td>
              <td class="px-4 py-3 text-sm text-slate-600 text-center">{{ item.count }}</td>
              <td class="px-4 py-3 text-sm font-bold text-red-600 font-mono text-right">{{ formatRupiah(item.total) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>

import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import api from '@/api/axiosInstance'
import BaseButton from '@/components/common/BaseButton.vue'
import { 
  Search, 
  DollarSign, 
  Barcode, 
  ArrowUpRight, 
  TrendingUp, 
  ClipboardList 
} from 'lucide-vue-next'

// State
const loading = ref(false)
const startDate = ref('')
const endDate = ref('')
const report = ref({
  total_revenue: 0,
  total_hpp: 0,
  gross_profit: 0,
  total_expense: 0,
  net_profit: 0,
  transaction_count: 0,
  expense_breakdown: []
})

// Set default dates (1st of current month to today)
const setDefaultDates = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  
  startDate.value = `${year}-${month}-01`
  endDate.value = `${year}-${month}-${day}`
}

// Fetch report
const fetchReport = async () => {
  loading.value = true
  try {
    const response = await api.get('/reports', {
      params: {
        start_date: startDate.value,
        end_date: endDate.value
      }
    })
    
    if (response.data && response.data.data) {
      report.value = response.data.data
    }
  } catch (error) {
    console.error('Error fetching report:', error)
    alert('Gagal mengambil data laporan')
  } finally {
    loading.value = false
  }
}

// Format currency
const formatRupiah = (value) => {
  const num = parseFloat(value) || 0
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(num)
}

// Lifecycle
onMounted(() => {
  setDefaultDates()
  fetchReport()
})
</script>

<style scoped>
/* Custom styles if needed */
</style>

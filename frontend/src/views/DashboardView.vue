<template>
  <AppLayout>
    <!-- Header Section -->
    <div v-if="summary" class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">Dashboard Owner</h1>
        <p class="text-slate-500 text-sm">Ringkasan performa toko hari ini</p>
      </div>
      <div class="flex items-center gap-2">
        <BaseButton 
          variant="ghost" 
          size="icon" 
          @click="fetchData" 
          :disabled="isLoading" 
          title="Refresh Data"
        >
          <RefreshCcw class="w-5 h-5" :class="{'animate-spin': isLoading}" />
        </BaseButton>
      </div>
    </div>

    <!-- Error Alert -->
    <!-- Error Alert -->
    <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg flex items-center gap-3">
      <TriangleAlert class="w-5 h-5" />
      <p class="text-sm font-medium">{{ error }}</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatCard 
        title="Omzet Hari Ini" 
        :value="formatRupiah(summary.total_revenue)" 
        color="blue"
        :loading="isLoading"
      >
        <template #icon>
          <DollarSign class="w-5 h-5" />
        </template>
      </StatCard>

      <StatCard 
        title="Pengeluaran" 
        :value="formatRupiah(summary.total_expense)" 
        color="red"
        :loading="isLoading"
      >
        <template #icon>
          <Receipt class="w-5 h-5" />
        </template>
      </StatCard>

      <StatCard 
        title="Laba Bersih" 
        :value="formatRupiah(summary.net_profit)" 
        color="green"
        :loading="isLoading"
      >
        <template #icon>
          <Wallet class="w-5 h-5" />
        </template>
      </StatCard>

      <StatCard 
        title="Total Transaksi" 
        :value="summary.transaction_count + ' Struk'" 
        color="orange"
        :loading="isLoading"
      >
        <template #icon>
          <ShoppingCart class="w-5 h-5" />
        </template>
      </StatCard>
    </div>

    <!-- Charts & Tables Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Sales Chart -->
      <div class="lg:col-span-2 dashboard-card p-6">
        <div class="flex justify-between items-center mb-6">
          <div>
            <h3 class="text-lg font-bold text-slate-800">Analitik Penjualan</h3>
            <p class="text-xs text-slate-500 font-medium mt-1">Tren pendapatan 7 hari terakhir</p>
          </div>
        </div>
        <div class="h-80 w-full relative">
          <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-white/50 z-10">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
          <Line v-if="chartData.labels.length > 0" :data="chartData" :options="chartOptions" />
          <div v-else class="h-full flex items-center justify-center text-slate-400 text-sm">
            Belum ada data penjualan minggu ini
          </div>
        </div>
      </div>

      <!-- Top Products -->
      <div class="dashboard-card overflow-hidden flex flex-col">
        <div class="p-6 border-b border-gray-100 flex justify-between items-center">
          <h3 class="font-bold text-slate-800 flex items-center gap-2">
            <Trophy class="w-5 h-5 text-yellow-500" /> Produk Terlaris
          </h3>
        </div>
        
        <div class="flex-1 overflow-auto custom-scroll p-0">
          <table class="w-full text-left">
            <thead class="bg-slate-50 text-slate-500 text-xs uppercase font-bold tracking-wider border-b border-slate-100">
              <tr>
                <th class="px-4 py-3">Produk</th>
                <th class="px-4 py-3 text-center">Qty</th>
                <th class="px-4 py-3 text-right">Total</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-sm">
              <tr v-if="isLoading" v-for="i in 5" :key="i">
                <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded w-24 animate-pulse"></div></td>
                <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded w-8 mx-auto animate-pulse"></div></td>
                <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded w-16 ml-auto animate-pulse"></div></td>
              </tr>
              <tr v-else-if="topProducts.length === 0">
                <td colspan="3" class="px-4 py-8 text-center text-slate-400 text-xs">Belum ada data penjualan</td>
              </tr>
              <tr 
                v-else 
                v-for="product in topProducts" 
                :key="product.product_id"
                class="hover:bg-slate-50 transition"
              >
                <td class="px-4 py-3">
                  <p class="font-semibold text-slate-700 truncate max-w-[120px]" :title="product.product_name">
                    {{ product.product_name }}
                  </p>
                </td>
                <td class="px-4 py-3 text-center font-bold text-slate-600">
                  {{ product.total_qty }}
                </td>
                <td class="px-4 py-3 text-right font-bold text-blue-600">
                  {{ formatRupiahSimple(product.total_sales) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatCard from '@/components/dashboard/StatCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import api from '@/api/axiosInstance'
import {
  RefreshCcw,
  TriangleAlert,
  DollarSign,
  Receipt,
  Wallet,
  ShoppingCart,
  Trophy
} from 'lucide-vue-next'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line } from 'vue-chartjs'

// Register ChartJS Components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

// State
const isLoading = ref(false)
const error = ref(null)

const summary = ref({
  total_revenue: 0,
  total_expense: 0,
  net_profit: 0,
  transaction_count: 0
})

const topProducts = ref([])

// Chart State
const chartData = ref({
  labels: [],
  datasets: []
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: function(context) {
          let label = context.dataset.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed.y !== null) {
            label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(context.parsed.y);
          }
          return label;
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: { borderDash: [4, 4], color: '#e2e8f0', drawBorder: false },
      ticks: {
        callback: function(value) {
          if (value >= 1000000) return (value/1000000) + 'jt';
          if (value >= 1000) return (value/1000) + 'rb';
          return value;
        },
        color: '#94a3b8',
        font: { size: 10, weight: 500 }
      }
    },
    x: {
      grid: { display: false },
      ticks: { color: '#64748b', font: { size: 10, weight: 500 } }
    }
  },
  elements: {
    line: { tension: 0.3 },
    point: { radius: 4, hoverRadius: 6 }
  }
}

// Helpers
const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}

const formatRupiahSimple = (value) => {
  if (value >= 1000000) return (value / 1000000).toFixed(1) + 'jt'
  if (value >= 1000) return (value / 1000).toFixed(0) + 'rb'
  return value
}

// Data Fetching
const fetchData = async () => {
  isLoading.value = true
  error.value = null

  try {
    // 1. Get Summary
    const summaryRes = await api.get('/dashboard/summary')
    if (summaryRes.data.status === 'success') {
      summary.value = summaryRes.data.data
    }

    // 2. Get Sales Chart
    const chartRes = await api.get('/dashboard/chart?days=7')
    if (chartRes.data.status === 'success') {
      const data = chartRes.data.data.data
      
      chartData.value = {
        labels: data.map(d => {
          const date = new Date(d.date)
          return date.toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric' })
        }),
        datasets: [{
          label: 'Revenue',
          data: data.map(d => d.revenue),
          borderColor: '#2563eb',
          backgroundColor: (context) => {
            const ctx = context.chart.ctx;
            const gradient = ctx.createLinearGradient(0, 0, 0, 300);
            gradient.addColorStop(0, 'rgba(37, 99, 235, 0.2)');
            gradient.addColorStop(1, 'rgba(37, 99, 235, 0)');
            return gradient;
          },
          borderWidth: 2,
          pointBackgroundColor: '#fff',
          pointBorderColor: '#2563eb',
          fill: true
        }]
      }
    }

    // 3. Get Top Products
    const topRes = await api.get('/dashboard/top-products?limit=5')
    if (topRes.data.status === 'success') {
      topProducts.value = topRes.data.data.products
    }

  } catch (err) {
    console.error('Fetch dashboard error:', err)
    error.value = 'Gagal memuat data dashboard. Periksa koneksi internet.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard-card {
  background: white;
  border: 1px solid #E2E8F0;
  border-radius: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
}
</style>

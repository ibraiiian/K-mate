<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="mb-6 flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-slate-800">Riwayat Keuangan</h1>
          <p class="text-slate-500 text-sm">Monitor arus kas masuk dan keluar</p>
        </div>
        <BaseButton 
          variant="destructive"
          @click="openExpenseModal"
          class="shadow-lg shadow-red-500/20"
        >
          <MinusCircle class="w-4 h-4 mr-2" /> Catat Beban
        </BaseButton>
      </div>

    <!-- 1. SUMMARY CARDS (Quick Cashflow Overview) -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
      <div class="bg-white border border-gray-200 p-4 rounded-lg shadow-sm flex items-center justify-between">
        <div>
          <p class="text-xs text-slate-500 font-bold uppercase">Total Masuk (Income)</p>
          <p class="text-lg font-bold text-green-600 font-mono mt-1">+ {{ formatRupiah(totalIncome) }}</p>
        </div>
        <div class="bg-green-50 p-2 rounded-md text-green-600"><ArrowDownLeft class="w-6 h-6" /></div>
      </div>
      <div class="bg-white border border-gray-200 p-4 rounded-lg shadow-sm flex items-center justify-between">
        <div>
          <p class="text-xs text-slate-500 font-bold uppercase">Total Keluar (Expense)</p>
          <p class="text-lg font-bold text-red-600 font-mono mt-1">- {{ formatRupiah(totalExpense) }}</p>
        </div>
        <div class="bg-red-50 p-2 rounded-md text-red-600"><ArrowUpRight class="w-6 h-6" /></div>
      </div>
      <div class="bg-slate-800 text-white p-4 rounded-lg shadow-sm flex items-center justify-between border border-slate-900">
        <div>
          <p class="text-xs text-slate-400 font-bold uppercase">Saldo Bersih (Net)</p>
          <p class="text-lg font-bold text-white font-mono mt-1">{{ formatRupiah(netBalance) }}</p>
        </div>
        <div class="bg-slate-700 p-2 rounded-md text-blue-300"><Wallet class="w-6 h-6" /></div>
      </div>
    </div>

    <!-- 2. FILTERS & ACTIONS -->
    <div class="flex flex-col xl:flex-row justify-between items-start xl:items-center gap-4 mb-6">
      
      <!-- Tabs Switcher -->
      <div class="flex p-1 bg-gray-200 rounded-lg border border-gray-300 w-full md:w-auto">
        <button 
          @click="activeFilter = 'all'" 
          :class="['tab-btn flex-1 md:flex-none', activeFilter === 'all' ? 'active' : '']"
        >Semua</button>
        <button 
          @click="activeFilter = 'in'" 
          :class="['tab-btn flex-1 md:flex-none', activeFilter === 'in' ? 'active' : '']"
        >Pemasukan</button>
        <button 
          @click="activeFilter = 'out'" 
          :class="['tab-btn flex-1 md:flex-none', activeFilter === 'out' ? 'active' : '']"
        >Pengeluaran</button>
      </div>

      <!-- Search & Date -->
      <div class="flex gap-2 w-full md:w-auto">
        <div class="relative flex-1 md:w-72">
          <Search class="absolute left-3 top-2.5 text-slate-400 w-4 h-4" />
          <BaseInput 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari No. Ref / Keterangan..." 
            class="pl-10 w-full"
          />
        </div>
        <div class="flex items-center gap-2 border border-slate-200 rounded-md px-3 py-2 bg-white">
          <Calendar class="w-4 h-4 text-slate-400" />
          <input 
            v-model="filterDate"
            type="date" 
            class="text-sm border-none focus:ring-0 text-slate-600 p-0 w-28 bg-transparent"
          >
        </div>
        
        <!-- Refresh Button -->
        <BaseButton 
          variant="outline"
          size="icon"
          @click="fetchData" 
          :disabled="loading"
          title="Refresh Data"
        >
          <RefreshCcw class="w-4 h-4" :class="{'animate-spin': loading}" />
        </BaseButton>


      </div>
    </div>

    <!-- 3. DATA TABLE -->
    <div class="bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden">
      
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center p-10">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredData.length === 0" class="flex flex-col items-center justify-center p-10 text-slate-400">
        <Receipt class="w-12 h-12 mb-2 opacity-50" />
        <p>Tidak ada data transaksi</p>
      </div>

      <!-- Table -->
      <div v-else class="overflow-x-auto">
        <table class="w-full data-table">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>No. Ref</th>
              <th>Keterangan</th>
              <th>Kategori</th>
              <th class="text-center">Tipe</th>
              <th class="text-right">Nominal</th>
              <th>Status</th>
              <th class="text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr 
              v-for="item in filteredData" 
              :key="item.id"
              :class="[
                'hover:bg-slate-50 cursor-pointer transition-colors',
                item.type === 'out' ? 'bg-red-50/30' : ''
              ]"
              @click="openDetailModal(item)"
            >
              <td class="text-slate-500">{{ formatTime(item.date) }}</td>
              <td class="font-mono text-slate-600 text-xs">{{ item.ref_no }}</td>
              <td class="font-medium text-slate-800">{{ item.description }}</td>
              <td>
                <span :class="['tag', getCategoryClass(item.category)]">{{ item.category }}</span>
              </td>
              <td class="text-center">
                <span v-if="item.type === 'in'" class="badge badge-in text-xs flex items-center">
                  <ArrowDownLeft class="w-3 h-3 mr-1" /> Masuk
                </span>
                <span v-else class="badge badge-out text-xs flex items-center">
                  <ArrowUpRight class="w-3 h-3 mr-1" /> Keluar
                </span>
              </td>
              <td class="text-right font-bold font-mono" :class="item.type === 'in' ? 'text-green-600' : 'text-red-600'">
                {{ item.type === 'in' ? '+' : '-' }} {{ formatRupiahShort(item.amount) }}
              </td>
              <td>
                <span class="text-xs font-bold" :class="item.type === 'in' ? 'text-green-600' : 'text-slate-500'">
                  {{ item.status }}
                </span>
              </td>
              <td class="text-center">
                <button 
                  @click.stop="openDetailModal(item)"
                  class="text-slate-500 hover:text-blue-600 transition"
                >
                  <Info class="w-5 h-5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- Pagination -->
      <div v-if="filteredData.length > 0" class="px-6 py-4 border-t border-gray-200 flex justify-between items-center bg-gray-50">
        <p class="text-xs text-slate-500">Menampilkan <span class="font-bold text-slate-700">{{ filteredData.length }}</span> data</p>
        <div class="flex gap-1">
          <button class="px-3 py-1 border border-gray-300 rounded bg-white text-slate-600 hover:bg-gray-50 text-xs font-medium disabled:opacity-50" disabled>Sebelumnya</button>
          <button class="px-3 py-1 border border-gray-300 rounded bg-white text-slate-600 hover:bg-gray-50 text-xs font-medium">Selanjutnya</button>
        </div>
      </div>
    </div>

    <!-- DETAIL MODAL -->
    <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="closeModal"></div>
      
      <!-- Modal Content -->
      <div class="relative bg-white rounded-lg shadow-xl w-[600px] max-h-[90vh] overflow-hidden border border-gray-200 animate-fadeSlideIn">
        <!-- Header -->
        <div class="bg-slate-50 px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <div>
            <h3 class="font-bold text-slate-800 text-lg">{{ selectedItem?.ref_no }}</h3>
            <p class="text-xs text-slate-500 mt-0.5">Detail Transaksi</p>
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
              <p class="font-medium text-slate-800">{{ formatDate(selectedItem?.date) }}</p>
            </div>
            <div>
              <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Tipe</p>
              <p class="font-medium" :class="selectedItem?.type === 'in' ? 'text-green-600' : 'text-red-600'">
                {{ selectedItem?.type === 'in' ? 'Pemasukan' : 'Pengeluaran' }}
              </p>
            </div>
            <div>
              <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Kategori</p>
              <p class="font-medium text-slate-800">{{ selectedItem?.category }}</p>
            </div>
            <div>
              <p class="text-slate-500 text-xs uppercase tracking-wide mb-1">Status</p>
              <p class="font-medium text-slate-800">{{ selectedItem?.status }}</p>
            </div>
          </div>

          <!-- Items Table (for transactions/sales) -->
          <div v-if="selectedItem?.type === 'in' && selectedItem?.items?.length > 0" class="border border-gray-200 rounded-md overflow-hidden mb-6">
            <table class="w-full text-sm text-left">
              <thead class="bg-gray-50 text-gray-500 border-b border-gray-200">
                <tr>
                  <th class="px-4 py-2 font-medium">Barang</th>
                  <th class="px-4 py-2 font-medium text-center">Jml</th>
                  <th class="px-4 py-2 font-medium text-right">Harga</th>
                  <th class="px-4 py-2 font-medium text-right">Subtotal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="(lineItem, idx) in selectedItem.items" :key="idx">
                  <td class="px-4 py-3 text-slate-800">{{ lineItem.product_name || lineItem.name }}</td>
                  <td class="px-4 py-3 text-center text-slate-600">{{ lineItem.quantity || lineItem.qty }}</td>
                  <td class="px-4 py-3 text-right text-slate-600 font-mono text-xs">{{ formatRupiah(lineItem.unit_price || lineItem.price) }}</td>
                  <td class="px-4 py-3 text-right text-slate-800 font-mono font-medium text-xs">{{ formatRupiah(lineItem.subtotal || (lineItem.quantity * lineItem.unit_price)) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Description (for expenses) -->
          <div v-else-if="selectedItem?.type === 'out'" class="bg-slate-50 p-4 rounded-lg border border-gray-200 mb-6">
            <p class="text-slate-500 text-xs uppercase tracking-wide mb-2">Keterangan Pengeluaran</p>
            <p class="text-slate-800">{{ selectedItem?.description }}</p>
          </div>

          <!-- Totals -->
          <div class="flex justify-end">
            <div class="w-1/2 space-y-2">
              <div class="flex justify-between text-base font-bold text-slate-900 border-t border-gray-200 pt-2 mt-2">
                <span>TOTAL</span>
                <span :class="selectedItem?.type === 'in' ? 'text-green-600' : 'text-red-600'">
                  {{ selectedItem?.type === 'in' ? '+' : '-' }} {{ formatRupiah(selectedItem?.amount) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex justify-between items-center">
          <BaseButton 
            v-if="selectedItem?.type === 'in' && selectedItem?.items?.length > 0"
            @click="openReturnModal"
            variant="destructive"
            class="bg-red-100 text-red-700 hover:bg-red-200 border-none"
          >
            <CornerUpLeft class="w-4 h-4 mr-2" /> Ajukan Retur
          </BaseButton>
          <div v-else></div>
          <div class="flex gap-3">
            <BaseButton 
              v-if="selectedItem?.type === 'in'"
              @click="reprintReceipt"
              class="bg-slate-800 text-white hover:bg-slate-900"
            >
              <Printer class="w-4 h-4 mr-2" /> Cetak Ulang Struk
            </BaseButton>
            <BaseButton 
              @click="closeModal"
              variant="secondary"
            >
              Tutup
            </BaseButton>
          </div>
        </div>
      </div>
    </div>

    <!-- ADD EXPENSE MODAL -->
    <div v-if="showExpenseModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="closeExpenseModal"></div>
      
      <div class="relative bg-white rounded-lg shadow-xl w-[500px] overflow-hidden border border-gray-200 animate-fadeSlideIn">
        <div class="bg-slate-50 px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="font-bold text-slate-800 text-lg">Catat Pengeluaran</h3>
          <button @click="closeExpenseModal" class="text-slate-400 hover:text-red-600 transition">
            <X class="w-6 h-6" />
          </button>
        </div>

        <form @submit.prevent="submitExpense" class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Judul Pengeluaran</label>
            <BaseInput v-model="newExpense.title" required placeholder="Contoh: Beli Token Listrik" />
          </div>
          
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Kategori</label>
              <select v-model="newExpense.category" class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm focus:outline-none focus:border-blue-600">
                <option value="operational">Operasional</option>
                <option value="salary">Gaji Karyawan</option>
                <option value="restock">Restock Barang</option>
                <option value="maintenance">Perawatan</option>
                <option value="other">Lainnya</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Tanggal</label>
              <input v-model="newExpense.expense_date" type="date" required class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm focus:outline-none focus:border-blue-600">
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Nominal (Rp)</label>
            <BaseInput v-model.number="newExpense.amount" type="number" required min="1" class="font-mono" placeholder="0" />
          </div>

          <div>
             <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Keterangan (Opsional)</label>
             <textarea v-model="newExpense.description" rows="2" class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm focus:outline-none focus:border-blue-600" placeholder="Detail tambahan..."></textarea>
          </div>

          <div class="pt-4 flex justify-end gap-3">
             <BaseButton type="button" variant="secondary" @click="closeExpenseModal">Batal</BaseButton>
             <BaseButton type="submit" variant="destructive" :disabled="submitting" :loading="submitting">
               Simpan Pengeluaran
             </BaseButton>
          </div>
        </form>
      </div>
    </div>
    </div>

    <!-- RETURN MODAL -->
    <ReturnModal 
      v-if="showReturnModal" 
      :transaction="selectedTransactionForReturn"
      @close="showReturnModal = false"
      @success="handleReturnSuccess"
    />
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ReturnModal from '@/components/pos/ReturnModal.vue'
import api from '@/api/axiosInstance'

import { usePrint } from '@/composables/usePrint'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import { 
  ArrowDownLeft, 
  ArrowUpRight, 
  Wallet, 
  Search, 
  Calendar, 
  RefreshCcw, 
  MinusCircle, 
  Receipt, 
  Info, 
  X, 
  Printer, 
  CornerUpLeft 
} from 'lucide-vue-next'

const { printTransaction } = usePrint()

// State
const loading = ref(false)
const cashflowData = ref([])
const activeFilter = ref('all')
const searchQuery = ref('')
const filterDate = ref('')

// Modal State
const showDetailModal = ref(false)
const selectedItem = ref(null)

// Return Modal State
const showReturnModal = ref(false)
const selectedTransactionForReturn = ref(null)

// Expense Modal State
const showExpenseModal = ref(false)
const submitting = ref(false)
const newExpense = ref({
  title: '',
  category: 'operational',
  amount: '',
  expense_date: new Date().toISOString().split('T')[0],
  description: ''
})

// Fetch Data
const fetchData = async () => {
  loading.value = true
  try {
    // Fetch transactions (income) and expenses in parallel
    const [transactionsRes, expensesRes] = await Promise.all([
      api.get('/pos/transactions'),
      api.get('/expenses')
    ])

    const transactionsRaw = transactionsRes.data?.data || []
    const transactions = Array.isArray(transactionsRaw) ? transactionsRaw : (transactionsRaw.transactions || [])
    
    const expensesRaw = expensesRes.data?.data || []
    // Handle both array (non-paginated) and object (paginated) response formats
    const expenses = Array.isArray(expensesRaw) ? expensesRaw : (expensesRaw.expenses || [])

    // Format transactions as income (IN)
    const formattedTransactions = transactions.map(t => ({
      id: `txn-${t.id}`,
      date: t.transaction_date || t.created_at,
      ref_no: t.invoice_no,
      description: t.customer_name ? `${t.customer_name} (${t.payment_method?.toUpperCase()})` : `Penjualan (${t.payment_method?.toUpperCase()})`,
      category: 'Penjualan',
      type: 'in',
      amount: parseFloat(t.total_amount) || 0,
      status: 'Lunas',
      items: t.items || [],
      original: t
    }))

    // Format expenses as expense (OUT)
    const formattedExpenses = expenses.map(e => ({
      id: `exp-${e.id}`,
      date: e.created_at || e.expense_date, // Prefer created_at for real-time
      ref_no: `EXP-${String(e.id).padStart(6, '0')}`,
      description: e.title || e.description || 'Pengeluaran', // Fallback title
      category: e.category || 'Operasional',
      type: 'out',
      amount: parseFloat(e.amount) || 0,
      status: 'Selesai',
      items: [],
      original: e,
      full_description: e.description // Keep full description for detail
    }))

    // Merge and sort by date descending
    cashflowData.value = [...formattedTransactions, ...formattedExpenses].sort((a, b) => {
      // Sort by newest first
      return new Date(b.date) - new Date(a.date)
    })

  } catch (error) {
    console.error('Error fetching cashflow data:', error)
  } finally {
    loading.value = false
  }
}

// Computed: Filtered Data
const filteredData = computed(() => {
  let result = cashflowData.value

  // Filter by type
  if (activeFilter.value !== 'all') {
    result = result.filter(item => item.type === activeFilter.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(item => 
      item.ref_no.toLowerCase().includes(query) ||
      item.description.toLowerCase().includes(query)
    )
  }

  // Filter by date
  if (filterDate.value) {
    const selectedDate = new Date(filterDate.value).toDateString()
    result = result.filter(item => {
      const itemDate = new Date(item.date).toDateString()
      return itemDate === selectedDate
    })
  }

  return result
})

// Computed: Totals
const totalIncome = computed(() => {
  return cashflowData.value
    .filter(item => item.type === 'in')
    .reduce((sum, item) => sum + item.amount, 0)
})

const totalExpense = computed(() => {
  return cashflowData.value
    .filter(item => item.type === 'out')
    .reduce((sum, item) => sum + item.amount, 0)
})

const netBalance = computed(() => {
  return totalIncome.value - totalExpense.value
})

// Modal Functions
const openDetailModal = async (item) => {
  selectedItem.value = item
  showDetailModal.value = true
  
  // For sales transactions, fetch full details including items
  if (item.type === 'in' && item.original?.id) {
    try {
      const response = await api.get(`/pos/transactions/${item.original.id}`)
      const detail = response.data?.data
      if (detail && detail.items) {
        // Map items with required fields for return
        selectedItem.value = {
          ...item,
          items: detail.items.map((i, idx) => ({
            id: i.id || idx + 1,
            detail_id: i.id || idx + 1,
            product_unit_id: i.product_unit_id,
            product_name: i.product_name,
            name: i.product_name,
            unit_name: i.unit_name,
            quantity: i.qty,
            qty: i.qty,
            unit_price: parseFloat(i.price_at_moment) || 0,
            price: parseFloat(i.price_at_moment) || 0,
            price_at_moment: parseFloat(i.price_at_moment) || 0,
            subtotal: parseFloat(i.subtotal) || 0
          })),
          transaction_id: detail.id
        }
      }
    } catch (error) {
      console.error('Error fetching transaction details:', error)
    }
  }
}

const closeModal = () => {
  showDetailModal.value = false
  selectedItem.value = null
}

// Reprint Receipt Function
const reprintReceipt = () => {
  if (selectedItem.value?.type === 'in' && selectedItem.value?.original) {
    const txn = selectedItem.value.original
    const items = selectedItem.value.items || []
    
    // Format data for receipt
    const receiptData = {
      invoiceNo: txn.invoice_no,
      date: txn.transaction_date,
      cashierName: txn.cashier_name || 'Admin',
      customerName: selectedItem.value.description?.split(' (')[0] || 'Pelanggan Umum',
      paymentMethod: txn.payment_method,
      items: items.map(item => ({
        productName: item.product_name || item.name,
        qty: item.qty || item.quantity || 1,
        price: item.price_at_moment || item.price || item.unit_price || 0,
        unitName: item.unit_name || 'Pcs'
      })),
      totalItems: items.length,
      totalPrice: parseFloat(txn.total_amount) || 0,
      cashGiven: parseFloat(txn.total_amount) || 0,
      change: 0
    }
    
    printTransaction(receiptData)
  }
}

// Expense Modal Functions
const openExpenseModal = () => {
    showExpenseModal.value = true
    // Reset form
    newExpense.value = {
        title: '',
        category: 'operational',
        amount: '',
        expense_date: new Date().toISOString().split('T')[0],
        description: ''
    }
}

const closeExpenseModal = () => {
    showExpenseModal.value = false
}

// Return Modal Functions
const openReturnModal = () => {
  selectedTransactionForReturn.value = selectedItem.value
  showDetailModal.value = false
  showReturnModal.value = true
}

const handleReturnSuccess = () => {
  fetchData() // Refresh data after successful return
}

const submitExpense = async () => {
    submitting.value = true
    try {
        const payload = {
            title: newExpense.value.title,
            category: newExpense.value.category,
            amount: parseFloat(newExpense.value.amount),
            expense_date: newExpense.value.expense_date,
            description: newExpense.value.description
        }

        const response = await api.post('/expenses', payload)
        
        if (response.data.status === 'success') {
            alert('Pengeluaran berhasil dicatat!')
            closeExpenseModal()
            fetchData() // Refresh list
        }
    } catch (error) {
        console.error('Error saving expense:', error)
        alert(error.response?.data?.message || 'Gagal menyimpan pengeluaran')
    } finally {
        submitting.value = false
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

const formatRupiahShort = (value) => {
  return new Intl.NumberFormat('id-ID', {
    minimumFractionDigits: 0
  }).format(value || 0)
}

const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  
  // Check if date string contains time (has T or space separator with time)
  const hasTime = dateStr.includes('T') || (dateStr.includes(' ') && dateStr.split(' ')[1]?.includes(':'))
  
  if (hasTime) {
    return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
  } else {
    // For date-only strings, show the date instead
    return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { 
    day: 'numeric', 
    month: 'long', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getCategoryClass = (category) => {
  const lower = category?.toLowerCase() || ''
  if (lower.includes('penjualan') || lower.includes('sales')) return 'tag-sales'
  if (lower.includes('restock') || lower.includes('pembelian')) return 'tag-restock'
  return 'tag-expense'
}

// Lifecycle
onMounted(() => {
  fetchData()
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
  border-bottom: 1px solid #E2E8F0;
  font-size: 0.875rem;
  color: #1E293B;
}
.data-table tr:last-child td {
  border-bottom: none;
}
.data-table tr:hover td {
  background-color: #F8FAFC;
}

/* Badges */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}
.badge-in { background-color: #DCFCE7; color: #166534; border: 1px solid #BBF7D0; }
.badge-out { background-color: #FEE2E2; color: #991B1B; border: 1px solid #FECACA; }

/* Tags */
.tag {
  font-size: 0.7rem;
  padding: 1px 6px;
  border-radius: 3px;
  font-weight: 600;
  text-transform: uppercase;
  border: 1px solid transparent;
}
.tag-sales { background-color: #E0E7FF; color: #3730A3; border-color: #C7D2FE; }
.tag-expense { background-color: #FEF3C7; color: #92400E; border-color: #FDE68A; }
.tag-restock { background-color: #E0F2FE; color: #075985; border-color: #BAE6FD; }

/* Tab Buttons */
.tab-btn {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  color: #64748B;
  transition: all 0.15s;
}
.tab-btn.active {
  background-color: white;
  color: #0F52BA;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
  font-weight: 600;
}
.tab-btn:hover:not(.active) {
  color: #334155;
  background-color: rgba(255,255,255,0.5);
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

<template>
  <AppLayout>
  <!-- Main POS Container with fixed height to fit inside AppLayout -->
  <div class="h-[calc(100vh-140px)] w-full flex flex-col font-sans bg-[#F1F5F9] text-[#334155] rounded-xl overflow-hidden border border-slate-200">




    <!-- 2. MAIN LAYOUT -->
    <div class="flex-1 flex overflow-hidden">
        
        <!-- LEFT COLUMN: CATALOG (65%) -->
        <div class="w-[65%] flex flex-col border-r border-gray-200 bg-white">
            
            <!-- Search & Filter Bar -->
            <div class="p-4 border-b border-gray-200 bg-white shadow-sm z-10">
                <div class="flex gap-3">
                    <div class="relative flex-1">
                        <Search class="absolute left-3 top-3 text-slate-400 w-4 h-4" />
                        <input 
                            v-model="searchQuery"
                            @input="handleSearch"
                            type="text" 
                            placeholder="Cari nama barang, kode SKU, atau barcode..." 
                            class="w-full pl-10 pr-4 py-2 bg-slate-100 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-slate-900 focus:bg-white transition"
                        >
                    </div>
                </div>

                <!-- Categories (Visual Only for now) -->
                <div class="flex gap-2 mt-4 overflow-x-auto custom-scroll pb-1">
                    <button class="px-4 py-1.5 rounded-full bg-blue-600 text-white text-xs font-bold whitespace-nowrap shadow-sm">Semua</button>
                    <!-- Add functionality later if needed -->
                    <button class="px-4 py-1.5 rounded-full bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200 text-xs font-medium whitespace-nowrap transition">Semen & Pasir</button>
                    <button class="px-4 py-1.5 rounded-full bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200 text-xs font-medium whitespace-nowrap transition">Cat & Pelapis</button>
                </div>
            </div>

            <!-- Product Grid -->
            <div class="flex-1 overflow-y-auto p-4 bg-slate-50 custom-scroll">
                 <div v-if="productStore.loading" class="flex items-center justify-center h-64">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                 </div>

                 <div v-else-if="productStore.products.length === 0" class="flex flex-col items-center justify-center p-10 text-slate-400">
                    <PackageOpen class="w-10 h-10 mb-2" />
                    <p>Produk tidak ditemukan</p>
                 </div>

                <div v-else class="grid grid-cols-3 xl:grid-cols-4 gap-4">
                    <!-- Product Card -->
                    <div 
                        v-for="product in productStore.products"
                        :key="product.id"
                        class="product-card bg-white rounded-lg p-3 flex flex-col gap-2 relative group" 
                        @click="handleProductClick(product)"
                    >
                        <div class="aspect-square bg-slate-100 rounded-md overflow-hidden relative border border-slate-100">
                             <img 
                                v-if="product.image_path" 
                                :src="`http://localhost:3000${product.image_path}`" 
                                class="w-full h-full object-cover opacity-90 group-hover:opacity-100 group-hover:scale-105 transition duration-500"
                                @error="$event.target.src = '/placeholder.png'" 
                             >
                             <Package class="h-full w-full p-8 text-slate-300" v-else />
                             
                            <div class="absolute bottom-0 right-0 bg-slate-800 text-white text-[10px] px-1.5 py-0.5 rounded-tl-md font-bold">
                                Stok: {{ product.stock_qty }}
                            </div>
                        </div>
                        <div>
                            <h4 class="font-bold text-slate-700 text-sm leading-tight h-9 overflow-hidden line-clamp-2" :title="product.name">{{ product.name }}</h4>
                            <div class="flex justify-between items-end mt-1">
                                <p class="text-blue-700 font-bold text-sm">
                                    {{ formatRupiah(getLowestPrice(product.units)) }}
                                </p>
                                <span v-if="product.units && product.units.length > 1" class="text-[10px] bg-orange-100 text-orange-700 px-1 rounded font-bold border border-orange-200">Multi</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- RIGHT COLUMN: CART (35%) -->
        <div class="w-[35%] flex flex-col bg-white shadow-xl z-10">
            
            <!-- Customer Info -->
            <div class="h-16 px-4 border-b border-gray-200 flex items-center justify-between bg-white shrink-0">
                <div class="flex items-center gap-2 cursor-pointer hover:bg-slate-50 p-1.5 rounded transition w-full">
                    <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 shrink-0">
                        <User class="w-4 h-4" />
                    </div>
                    <div class="w-full">
                        <BaseInput 
                            v-model="cartStore.customerName" 
                            type="text" 
                            class="h-6 text-xs font-bold text-slate-700 w-full bg-transparent border-none p-0 focus-visible:ring-0 placeholder:text-slate-400"
                            placeholder="Pelanggan Umum (Klik utk edit)"
                        />
                        <p class="text-[10px] text-slate-400">Cash / Walk-in</p>
                    </div>
                </div>
            </div>

            <!-- Cart Items List -->
            <div class="flex-1 overflow-y-auto p-2 space-y-1 custom-scroll bg-slate-50/50">
                
                <!-- Empty State -->
                <div v-if="cartStore.cart.length === 0" class="h-full flex flex-col items-center justify-center text-slate-400">
                    <ShoppingCart class="w-16 h-16 mb-2 text-slate-300" />
                    <p class="text-sm font-medium">Keranjang Kosong</p>
                    <p class="text-xs">Pilih produk di sebelah kiri</p>
                </div>

                <!-- Items -->
                <div 
                    v-for="(item, index) in cartStore.cart" 
                    :key="`${item.productId}-${item.unitId}`"
                    class="bg-white p-3 rounded-lg border border-slate-200 shadow-sm flex justify-between items-center group hover:border-blue-300 transition-colors"
                >
                    <div class="flex-1 min-w-0 pr-2">
                        <h4 class="font-bold text-slate-700 text-sm truncate">{{ item.productName }}</h4>
                        <p class="text-xs text-slate-500">
                            {{ item.qty }} x {{ formatRupiah(item.price) }} / <span class="text-blue-600 font-bold">{{ item.unitName }}</span>
                        </p>
                    </div>
                    <div class="flex items-center gap-3">
                        <p class="font-bold text-slate-800 text-sm">{{ formatRupiah(item.price * item.qty) }}</p>
                        <div class="flex items-center border border-slate-300 rounded bg-slate-50">
                            <button 
                                @click="cartStore.updateQty(index, item.qty - 1)" 
                                :disabled="item.qty <= 1"
                                class="px-2 py-0.5 hover:bg-slate-200 text-slate-600 font-bold disabled:opacity-50"
                            >-</button>
                            <span class="px-1 text-xs font-bold w-6 text-center">{{ item.qty }}</span>
                            <button 
                                @click="cartStore.updateQty(index, item.qty + 1)"
                                class="px-2 py-0.5 hover:bg-slate-200 text-blue-600 font-bold border-l border-slate-300"
                            >+</button>
                        </div>
                        <button @click="cartStore.removeFromCart(index)" class="text-slate-300 hover:text-red-500 transition">
                            <Trash2 class="w-4 h-4" />
                        </button>
                    </div>
                </div>
            </div>

            <!-- Calculation Area -->
            <div class="bg-white border-t border-gray-200 p-4 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)] z-20">
                <!-- Subtotals -->
                <div class="space-y-2 mb-4 text-sm">
                    <div class="flex justify-between text-slate-500">
                        <span>Total Item</span>
                        <span class="font-medium text-slate-700">{{ cartStore.totalItems }} Items</span>
                    </div>
                    <!-- Assuming tax is included or simplified for now, as per backend -->
                    <div class="flex justify-between items-center pt-1 border-t border-dashed border-gray-200 mt-2">
                        <span class="text-base font-bold text-slate-800">Total Tagihan</span>
                        <span class="text-2xl font-bold text-blue-700">{{ formatRupiah(cartStore.totalPrice) }}</span>
                    </div>
                </div>

                <!-- Action Buttons -->
                <div class="grid grid-cols-3 gap-3">
                    <BaseButton 
                        variant="destructive"
                        class="col-span-1"
                        @click="confirmClearCart"
                    >
                        <Trash2 class="w-4 h-4 mr-2" />
                        Batal
                    </BaseButton>
                    <BaseButton 
                        variant="default"
                        class="col-span-2 bg-blue-700 hover:bg-blue-800"
                        :disabled="cartStore.cart.length === 0"
                        @click="openCheckout" 
                    >
                        Bayar Sekarang <ArrowRight class="w-4 h-4 ml-2" />
                    </BaseButton>
                </div>
            </div>
        </div>
    </div>


  </div>

    <!-- 3. MODALS (Dialogs) -->

    <!-- Unit Selector Component -->
    <UnitSelector 
      :show="showUnitSelector"
      :product="selectedProduct"
      @close="showUnitSelector = false"
      @select="handleUnitSelect"
    />

    <!-- Checkout Modal (Custom Implementation) -->
    <div v-if="showCheckoutModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="closeCheckout"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-white rounded-xl shadow-2xl w-[500px] overflow-hidden transform scale-100 transition-all">
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-white">
                <h3 class="font-bold text-lg text-slate-800">Pembayaran</h3>
                <button @click="closeCheckout" class="text-slate-400 hover:text-slate-600"><X class="w-5 h-5" /></button>
            </div>
            
            <div class="p-6">
                <div class="text-center mb-6">
                    <p class="text-sm text-slate-500">Total yang harus dibayar</p>
                    <h2 class="text-4xl font-bold text-slate-800 mt-1">{{ formatRupiah(cartStore.totalPrice) }}</h2>
                </div>

                <!-- Payment Methods -->
                <div class="grid grid-cols-3 gap-3 mb-6">
                    <button 
                        @click="paymentMethod = 'cash'"
                        class="py-3 border rounded-lg font-bold text-sm flex flex-col items-center justify-center transition-all"
                        :class="paymentMethod === 'cash' ? 'border-blue-600 bg-blue-50 text-blue-700 ring-2 ring-blue-600 ring-offset-1' : 'border-gray-200 text-slate-600 hover:bg-slate-50'"
                    >
                        <Banknote class="w-6 h-6 mb-1" /> Tunai
                    </button>
                    <button 
                        @click="paymentMethod = 'qris'"
                        class="py-3 border rounded-lg font-bold text-sm flex flex-col items-center justify-center transition-all"
                        :class="paymentMethod === 'qris' ? 'border-blue-600 bg-blue-50 text-blue-700 ring-2 ring-blue-600 ring-offset-1' : 'border-gray-200 text-slate-600 hover:bg-slate-50'"
                    >
                        <QrCode class="w-6 h-6 mb-1" /> QRIS
                    </button>
                    <button 
                        @click="paymentMethod = 'transfer'"
                        class="py-3 border rounded-lg font-bold text-sm flex flex-col items-center justify-center transition-all"
                        :class="paymentMethod === 'transfer' ? 'border-blue-600 bg-blue-50 text-blue-700 ring-2 ring-blue-600 ring-offset-1' : 'border-gray-200 text-slate-600 hover:bg-slate-50'"
                    >
                        <CreditCard class="w-6 h-6 mb-1" /> Transfer
                    </button>
                </div>

                <!-- Cash Input Area -->
                <div class="space-y-4">
                    <div>
                        <label class="text-xs font-bold text-slate-500 uppercase">Uang Diterima</label>
                        <input 
                            v-model.number="cashGiven"
                            type="number" 
                            class="w-full mt-1 p-3 border border-gray-300 rounded-lg text-lg font-bold text-slate-800 focus:outline-none focus:border-blue-600" 
                            placeholder="Rp 0"
                        >
                    </div>
                    
                    <!-- Quick Cash Buttons -->
                    <div class="flex gap-2">
                        <button @click="cashGiven = 50000" class="px-3 py-1 bg-slate-100 text-slate-600 text-xs font-bold rounded hover:bg-slate-200">Rp 50k</button>
                        <button @click="cashGiven = 100000" class="px-3 py-1 bg-slate-100 text-slate-600 text-xs font-bold rounded hover:bg-slate-200">Rp 100k</button>
                        <button @click="cashGiven = cartStore.totalPrice" class="px-3 py-1 bg-blue-100 text-blue-700 text-xs font-bold rounded hover:bg-blue-200">Uang Pas</button>
                    </div>

                    <!-- Change Display -->
                    <div class="p-4 bg-slate-50 rounded-lg flex justify-between items-center border border-slate-200">
                        <span class="font-bold text-slate-600">Kembalian</span>
                        <span 
                            class="font-bold text-xl"
                            :class="change >= 0 ? 'text-green-600' : 'text-red-500'"
                        >
                            {{ change >= 0 ? formatRupiah(change) : 'Kurang ' + formatRupiah(Math.abs(change)) }}
                        </span>
                    </div>
                </div>

                <BaseButton 
                    @click="processPayment"
                    :disabled="cartStore.loading || (paymentMethod === 'cash' && change < 0)"
                    :loading="cartStore.loading"
                    class="w-full mt-6 h-12 bg-green-600 hover:bg-green-700 text-white font-bold text-base"
                >
                    <Printer class="w-5 h-5 mr-2" /> Bayar & Cetak Struk
                </BaseButton>
            </div>
        </div>
    </div>

    <!-- Receipt Modal -->
    <div v-if="showReceiptModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-slate-900/70 backdrop-blur-sm" @click="closeReceipt"></div>
        
        <!-- Receipt Content -->
        <div class="relative bg-white rounded-xl shadow-2xl w-[380px] max-h-[90vh] overflow-hidden">
            <!-- Header -->
            <div class="bg-slate-800 px-6 py-4 text-white flex justify-between items-center">
                <div>
                    <h3 class="font-bold text-lg">Struk Pembayaran</h3>
                    <p class="text-slate-400 text-xs">{{ receiptData?.invoiceNo }}</p>
                </div>
                <button @click="closeReceipt" class="text-slate-400 hover:text-white transition">
                    <X class="w-5 h-5" />
                </button>
            </div>

            <!-- Receipt Body -->
            <div class="p-6 overflow-y-auto max-h-[60vh] custom-scroll">
                <!-- Store Info -->
                <div class="text-center border-b border-dashed border-gray-300 pb-4 mb-4">
                    <h2 class="font-bold text-xl text-slate-800">TB Karib</h2>
                    <p class="text-xs text-slate-500">Jl. Mt Haryono No. 41, Kota Cirebon, Sindang Laut</p>
                    <p class="text-xs text-slate-500">Telp: 021-1234567</p>
                </div>

                <!-- Transaction Info -->
                <div class="text-xs text-slate-600 space-y-1 mb-4 border-b border-dashed border-gray-300 pb-4">
                    <div class="flex justify-between">
                        <span>Tanggal</span>
                        <span class="font-medium">{{ receiptData?.date }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Kasir</span>
                        <span class="font-medium">Admin</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Pelanggan</span>
                        <span class="font-medium">{{ receiptData?.customerName }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Metode</span>
                        <span class="font-medium uppercase">{{ receiptData?.paymentMethod }}</span>
                    </div>
                </div>

                <!-- Items -->
                <div class="space-y-2 mb-4 border-b border-dashed border-gray-300 pb-4">
                    <div 
                        v-for="(item, index) in receiptData?.items" 
                        :key="index"
                        class="text-xs"
                    >
                        <div class="flex justify-between font-medium text-slate-700">
                            <span class="truncate max-w-[200px]">{{ item.productName }}</span>
                        </div>
                        <div class="flex justify-between text-slate-500 pl-2">
                            <span>{{ item.qty }} x {{ formatRupiah(item.price) }} ({{ item.unitName }})</span>
                            <span class="font-medium text-slate-700">{{ formatRupiah(item.price * item.qty) }}</span>
                        </div>
                    </div>
                </div>

                <!-- Totals -->
                <div class="space-y-2 text-sm">
                    <div class="flex justify-between text-slate-600">
                        <span>Total Item</span>
                        <span class="font-medium">{{ receiptData?.totalItems }} pcs</span>
                    </div>
                    <div class="flex justify-between font-bold text-lg text-slate-800 border-t border-gray-200 pt-2">
                        <span>TOTAL</span>
                        <span class="text-blue-700">{{ formatRupiah(receiptData?.totalPrice) }}</span>
                    </div>
                    <div class="flex justify-between text-slate-600">
                        <span>Tunai</span>
                        <span class="font-medium">{{ formatRupiah(receiptData?.cashGiven) }}</span>
                    </div>
                    <div class="flex justify-between text-green-600 font-bold">
                        <span>Kembalian</span>
                        <span>{{ formatRupiah(receiptData?.change) }}</span>
                    </div>
                </div>

                <!-- Footer Message -->
                <div class="mt-6 pt-4 border-t border-dashed border-gray-300 text-center">
                    <p class="text-xs text-slate-500">Terima kasih telah berbelanja!</p>
                    <p class="text-xs text-slate-400 mt-1">Barang yang sudah dibeli tidak dapat dikembalikan.</p>
                </div>
            </div>

            <!-- Actions -->
            <div class="p-4 border-t border-gray-200 bg-gray-50 flex gap-3">
                <button 
                    @click="printReceipt"
                    class="flex-1 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold text-sm flex items-center justify-center gap-2 transition"
                >
                    <Printer class="w-5 h-5" /> Cetak
                </button>
                <button 
                    @click="closeReceipt"
                    class="flex-1 py-3 bg-slate-200 hover:bg-slate-300 text-slate-700 rounded-lg font-bold text-sm flex items-center justify-center gap-2 transition"
                >
                    <Check class="w-5 h-5" /> Selesai
                </button>
            </div>
        </div>
    </div>

  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useProductStore } from '@/stores/productStore'
import { useCartStore } from '@/stores/cartStore'
import UnitSelector from '@/components/pos/UnitSelector.vue'
import { usePrint } from '@/composables/usePrint'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import { 
  Search, 
  Package, 
  PackageOpen, 
  ShoppingCart, 
  Trash2, 
  User, 
  ArrowRight,
  X,
  CreditCard,
  Banknote,
  QrCode,
  Printer,
  Check
} from 'lucide-vue-next'

const productStore = useProductStore()
const cartStore = useCartStore()
const { printTransaction } = usePrint()

// State
const searchQuery = ref('')
const currentTime = ref('00:00')
const showUnitSelector = ref(false)
const selectedProduct = ref(null)

// Checkout Modal State
const showCheckoutModal = ref(false)
const paymentMethod = ref('cash')
const cashGiven = ref(0) // Start with 0

// Receipt Modal State
const showReceiptModal = ref(false)
const receiptData = ref(null)

// Clock Logic
let clockInterval
const updateClock = () => {
    const now = new Date()
    currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

// Search Logic
let searchTimeout
const handleSearch = () => {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        productStore.fetchProducts(1, 10, searchQuery.value)
    }, 500)
}

// Product Interaction
const handleProductClick = (product) => {
    // Logic: If only 1 unit and it's the base unit or just one unit available, add immediately? 
    // HTML said: "if product.units.length === 1 -> addToCart".
    // Keeping it simple with UnitSelector for consistency, or implement HTML logic:
    if (product.units && product.units.length === 1) {
        cartStore.addToCart(product, product.units[0], 1)
    } else {
        selectedProduct.value = product
        showUnitSelector.value = true
    }
}

const handleUnitSelect = (unit) => {
    cartStore.addToCart(selectedProduct.value, unit, 1)
    showUnitSelector.value = false
    selectedProduct.value = null
}

// Cart Actions
const confirmClearCart = () => {
    if (confirm('Kosongkan keranjang?')) {
        cartStore.clearCart()
    }
}

// Checkout Logic
const openCheckout = () => {
    paymentMethod.value = 'cash'
    cashGiven.value = 0
    showCheckoutModal.value = true
}

const closeCheckout = () => {
    showCheckoutModal.value = false
}

const change = computed(() => {
    return cashGiven.value - cartStore.totalPrice
})

const processPayment = async () => {
    cartStore.paymentMethod = paymentMethod.value
    
    // Store cash/change info for receipt
    const currentCashGiven = cashGiven.value
    const currentChange = change.value
    const currentTotal = cartStore.totalPrice

    const success = await cartStore.checkout()
    if (success) {
        // Update receipt data with cash info
        if (cartStore.lastTransaction) {
            cartStore.lastTransaction.cashGiven = currentCashGiven
            cartStore.lastTransaction.change = currentChange
        }
        
        // Prepare receipt data
        receiptData.value = cartStore.lastTransaction
        
        // Close checkout, show receipt
        closeCheckout()
        showReceiptModal.value = true
    } else {
        alert(cartStore.error)
    }
}

// Receipt Functions
const closeReceipt = () => {
    showReceiptModal.value = false
    receiptData.value = null
}

const printReceipt = () => {
    // Use thermal print composable
    if (receiptData.value) {
        printTransaction(receiptData.value)
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

const getLowestPrice = (units) => {
    if (!units || units.length === 0) return 0
    return Math.min(...units.map(u => u.sell_price))
}

// Lifecycle
onMounted(() => {
    productStore.fetchProducts()
    updateClock()
    clockInterval = setInterval(updateClock, 1000)
})

onUnmounted(() => {
    clearInterval(clockInterval)
})
</script>

<style scoped>
/* Custom Scrollbar from HTML */
.custom-scroll::-webkit-scrollbar { width: 6px; height: 6px; }
.custom-scroll::-webkit-scrollbar-thumb { background: #CBD5E1; border-radius: 4px; }
.custom-scroll:hover::-webkit-scrollbar-thumb { background: #94a3b8; }

.product-card { transition: all 0.15s ease-out; cursor: pointer; border: 1px solid #E2E8F0; }
.product-card:active { transform: scale(0.98); border-color: #0F52BA; }
</style>

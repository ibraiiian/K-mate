<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 class="text-2xl font-bold text-slate-800">Manajemen Stok</h1>
          <p class="text-sm text-slate-500">Kelola daftar produk, harga, dan stok barang.</p>
        </div>
        <BaseButton 
          @click="openAddModal" 
          class="bg-blue-600 hover:bg-blue-700"
        >
          <Plus class="w-4 h-4 mr-2" />
          Tambah Produk
        </BaseButton>
      </div>

      <!-- Filters -->
      <div class="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col sm:flex-row gap-4">
        <div class="relative flex-1">
          <Search class="absolute left-3 top-2.5 text-gray-400 w-4 h-4" />
          <BaseInput 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari nama produk atau kode..." 
            class="pl-10"
            @input="handleSearch"
          />
        </div>
        <div class="flex gap-2">
           <BaseButton variant="outline" class="text-slate-600">
            <Filter class="w-4 h-4 mr-2" />
            Filter
          </BaseButton>
        </div>
      </div>

      <!-- Product List Table -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead class="bg-slate-50 text-slate-500 text-xs uppercase font-bold tracking-wider">
              <tr>
                <th class="px-6 py-4">Produk</th>
                <th class="px-6 py-4">Kode</th>
                <th class="px-6 py-4 text-center">Stok</th>
                <th class="px-6 py-4 text-right">Harga Jual</th>
                <th class="px-6 py-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 text-sm">
              <tr v-if="productStore.loading" class="animate-pulse">
                <td colspan="5" class="px-6 py-8 text-center text-slate-400">Memuat data...</td>
              </tr>
              <tr v-else-if="productStore.products.length === 0">
                <td colspan="5" class="px-6 py-8 text-center text-slate-400">Belum ada produk</td>
              </tr>
              
              <tr 
                v-else 
                v-for="product in productStore.products" 
                :key="product.id" 
                class="hover:bg-slate-50 transition group"
              >
                <!-- Product Name & Image -->
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="h-10 w-10 rounded-lg bg-gray-100 border border-gray-200 flex items-center justify-center overflow-hidden shrink-0">
                      <img 
                        v-if="product.image_path" 
                        :src="getImageUrl(product.image_path)" 
                        class="h-full w-full object-cover"
                        alt="Product Image"
                        @error="$event.target.src = '/placeholder.png'" 
                      >
                      <i v-else class="text-gray-400"><Package class="w-5 h-5" /></i>
                    </div>
                    <div>
                      <p class="font-bold text-slate-700">{{ product.name }}</p>
                      <p class="text-xs text-slate-500">{{ product.base_unit }}</p>
                    </div>
                  </div>
                </td>

                <!-- Code -->
                <td class="px-6 py-4">
                  <span class="font-mono text-slate-500 bg-slate-100 px-2 py-0.5 rounded text-xs">
                    {{ product.product_code || '-' }}
                  </span>
                </td>

                <!-- Stock -->
                <td class="px-6 py-4 text-center">
                  <span 
                    class="font-bold px-2.5 py-1 rounded-full text-xs"
                    :class="getStatusColor(product.stock_qty)"
                  >
                    {{ product.stock_qty }}
                  </span>
                </td>

                <!-- Price -->
                <td class="px-6 py-4 text-right font-medium text-slate-700">
                  {{ formatRupiah(product.units?.[0]?.sell_price) }}
                </td>

                <!-- Action -->
                <td class="px-6 py-4 text-center">
                  <button @click="openEditModal(product)" class="text-slate-400 hover:text-blue-600 transition p-1">
                    <Pencil class="w-4 h-4" />
                  </button>
                  <button @click="deleteProduct(product.id)" class="text-slate-400 hover:text-red-500 transition p-1 ml-2">
                    <Trash2 class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="px-6 py-4 border-t border-gray-100 flex items-center justify-between">
           <BaseButton 
             variant="ghost"
             :disabled="productStore.currentPage <= 1"
             @click="changePage(productStore.currentPage - 1)"
             class="text-slate-500 hover:text-blue-600"
            >
             <ChevronLeft class="w-4 h-4 mr-1" /> Previous
           </BaseButton>
           <span class="text-sm text-slate-500">Page {{ productStore.currentPage }}</span>
           <BaseButton 
             variant="ghost"
             :disabled="productStore.products.length < productStore.limit"
             @click="changePage(productStore.currentPage + 1)"
             class="text-slate-500 hover:text-blue-600"
            >
             Next <ChevronRight class="w-4 h-4 ml-1" />
           </BaseButton>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <ProductFormModal 
        :isOpen="showModal"
        :loading="productStore.loading"
        :product="selectedProduct"
        @close="closeModal"
        @submit="handleSubmit"
    />
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProductFormModal from '@/components/inventory/ProductFormModal.vue'
import { useProductStore } from '@/stores/productStore'
import { getImageUrl } from '@/utils/imageUrl'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import { 
  Plus, 
  Search, 
  Filter, 
  Package, 
  Pencil, 
  Trash2, 
  ChevronLeft, 
  ChevronRight 
} from 'lucide-vue-next'

const productStore = useProductStore()
const showModal = ref(false)
const selectedProduct = ref(null)
const searchQuery = ref('')

// Lifecycle
onMounted(() => {
    productStore.fetchProducts()
})

// Modal Handlers
const openAddModal = () => {
    selectedProduct.value = null
    showModal.value = true
}

const openEditModal = (product) => {
    selectedProduct.value = { ...product }
    showModal.value = true
}

const closeModal = () => {
    showModal.value = false
    selectedProduct.value = null
}

// Handlers
const getStatusColor = (qty) => {
    if (qty <= 5) return 'bg-red-50 text-red-600'
    if (qty <= 20) return 'bg-orange-50 text-orange-600'
    return 'bg-green-50 text-green-600'
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}

const handleSubmit = async (data, imageFile) => {
    let res
    
    if (selectedProduct.value?.id) {
        // Edit mode - pass imageFile too
        res = await productStore.updateProduct(data.id, data, imageFile)
    } else {
        // Add mode
        res = await productStore.addProduct(data)
    }
    
    if (res?.success) {
        closeModal()
    } else {
        alert(res?.message || 'Terjadi kesalahan')
    }
}

const deleteProduct = async (id) => {
    await productStore.deleteProduct(id)
}

// Debounce Search
let timeout
const handleSearch = () => {
    clearTimeout(timeout)
    timeout = setTimeout(() => {
        productStore.fetchProducts(1, 10, searchQuery.value)
    }, 500)
}

const changePage = (page) => {
    productStore.fetchProducts(page, 10, searchQuery.value)
}
</script>

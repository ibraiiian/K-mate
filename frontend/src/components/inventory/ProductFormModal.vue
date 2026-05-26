<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      
      <!-- Overlay -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeModal"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <!-- Modal Panel -->
      <div class="relative inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg w-full">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div class="sm:flex sm:items-start">
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
              <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                {{ isEditMode ? 'Edit Produk' : 'Tambah Produk Baru' }}
              </h3>
              
              <form @submit.prevent="handleSubmit" class="mt-4 space-y-4">
                
                <!-- Image Upload -->
                <div>
                    <label class="block text-sm font-medium text-gray-700">Foto Produk</label>
                    <div class="mt-1 flex items-center gap-4">
                        <div v-if="imagePreview" class="w-16 h-16 rounded-lg overflow-hidden border border-gray-200">
                            <img :src="imagePreview" class="w-full h-full object-cover">
                        </div>
                        <div v-else class="w-16 h-16 rounded-lg bg-gray-100 flex items-center justify-center text-gray-400">
                            <Image class="w-8 h-8" />
                        </div>
                        <label class="cursor-pointer bg-white py-2 px-3 border border-gray-300 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            Upload
                            <input type="file" class="hidden" @change="handleFileUpload" accept="image/*">
                        </label>
                    </div>
                </div>

                <!-- Product Name -->
                <div>
                  <label class="block text-sm font-medium text-gray-700">Nama Produk</label>
                  <input v-model="form.name" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                </div>

                <!-- Base Unit -->
                <div>
                  <label class="block text-sm font-medium text-gray-700">Satuan Dasar (Contoh: Pcs, Kg)</label>
                  <input v-model="form.base_unit" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                </div>

                <!-- Stock Qty -->
                <div>
                  <label class="block text-sm font-medium text-gray-700">Stok</label>
                  <input v-model="form.stock_qty" type="number" min="0" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                </div>

                <div class="border-t pt-4">
                    <h4 class="text-sm font-bold text-gray-700 mb-2">Harga & Konversi Satuan</h4>
                    <div class="bg-gray-50 p-3 rounded-lg space-y-3">
                        <div>
                            <label class="block text-xs font-medium text-gray-500">Nama Satuan (Contoh: Pcs)</label>
                            <input v-model="unitForm.unit_name" type="text" required placeholder="Pcs" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-1.5 px-3 text-sm">
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                             <div>
                                <label class="block text-xs font-medium text-gray-500">Harga Beli</label>
                                <input v-model="unitForm.buy_price" type="number" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-1.5 px-3 text-sm">
                            </div>
                             <div>
                                <label class="block text-xs font-medium text-gray-500">Harga Jual</label>
                                <input v-model="unitForm.sell_price" type="number" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-1.5 px-3 text-sm">
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="!isEditMode" class="bg-blue-50 p-3 rounded text-xs text-blue-700">
                    <p>Satuan dasar otomatis dianggap memiliki konversi 1. Anda dapat menambahkan satuan lain nanti.</p>
                </div>

                <div class="flex justify-end gap-3 mt-6">
                    <button type="button" @click="closeModal" class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-gray-700 bg-gray-100 hover:bg-gray-200 focus:outline-none">
                        Batal
                    </button>
                    <button type="submit" :disabled="loading" class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none">
                        {{ loading ? 'Menyimpan...' : (isEditMode ? 'Update Produk' : 'Simpan Produk') }}
                    </button>
                </div>

              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { Image } from 'lucide-vue-next'
import { getImageUrl } from '@/utils/imageUrl'

const props = defineProps({
  isOpen: Boolean,
  loading: Boolean,
  product: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['close', 'submit'])

const isEditMode = computed(() => props.product && props.product.id)

const form = reactive({
    name: '',
    base_unit: '',
    stock_qty: 0,
})

const unitForm = reactive({
    unit_name: '',
    conversion_factor: 1,
    buy_price: 0,
    sell_price: 0
})

const selectedFile = ref(null)
const imagePreview = ref(null)

// Reset function defined BEFORE watcher
const resetForm = () => {
    form.name = ''
    form.base_unit = ''
    form.stock_qty = 0
    unitForm.unit_name = ''
    unitForm.buy_price = 0
    unitForm.sell_price = 0
    selectedFile.value = null
    imagePreview.value = null
}

// Watch for product change to prefill form
watch(() => props.product, (newProduct) => {
    if (newProduct && newProduct.id) {
        form.name = newProduct.name || ''
        form.base_unit = newProduct.base_unit || ''
        form.stock_qty = newProduct.stock_qty || 0
        
        // Prefill unit if exists
        if (newProduct.units && newProduct.units.length > 0) {
            const unit = newProduct.units[0]
            unitForm.unit_name = unit.unit_name || ''
            unitForm.buy_price = unit.buy_price || 0
            unitForm.sell_price = unit.sell_price || 0
        }
        
        // Prefill image
        if (newProduct.image_path) {
            imagePreview.value = getImageUrl(newProduct.image_path)
        } else {
            imagePreview.value = null
        }
    } else {
        resetForm()
    }
}, { immediate: true })

const closeModal = () => {
    resetForm()
    emit('close')
}

const handleFileUpload = (event) => {
    const file = event.target.files[0]
    if (file) {
        selectedFile.value = file
        imagePreview.value = URL.createObjectURL(file)
    }
}

const handleSubmit = () => {
    if (isEditMode.value) {
        // Edit Mode - Send product data with ID and units
        const updateData = {
            id: props.product.id,
            name: form.name,
            base_unit: form.base_unit,
            stock_qty: parseInt(form.stock_qty),
            units: [
                {
                    id: props.product.units?.[0]?.id, // Include unit ID for update
                    unit_name: unitForm.unit_name,
                    conversion_factor: 1,
                    buy_price: parseFloat(unitForm.buy_price),
                    sell_price: parseFloat(unitForm.sell_price)
                }
            ]
        }
        emit('submit', updateData, selectedFile.value)
    } else {
        // Add Mode - Use FormData
        const productData = {
            name: form.name,
            base_unit: form.base_unit,
            stock_qty: parseInt(form.stock_qty),
            min_stock_alert: 10,
            units: [
                {
                    unit_name: unitForm.unit_name,
                    conversion_factor: 1,
                    buy_price: parseFloat(unitForm.buy_price),
                    sell_price: parseFloat(unitForm.sell_price)
                }
            ]
        }

        const formData = new FormData()
        if (selectedFile.value) {
            formData.append('image', selectedFile.value)
        }
        formData.append('data', JSON.stringify(productData))

        emit('submit', formData)
    }
}
</script>

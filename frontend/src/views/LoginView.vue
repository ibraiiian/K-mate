<template>
  <div class="bg-gray-50 text-gray-800 min-h-screen flex w-full font-sans">
    
    <!-- BAGIAN KIRI: Slider Visual Branding (Hidden on Mobile) -->
    <div 
      class="hidden lg:flex lg:w-1/2 relative overflow-hidden flex-col justify-between p-16 text-white transition-colors duration-1000" 
      :class="currentSlideData.bgColor"
    >
        
        <!-- Background Decorations -->
        <div class="absolute top-0 right-0 w-full h-full opacity-10 pointer-events-none">
            <!-- Blueprint Grid Pattern -->
            <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
                <defs>
                    <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
                        <path d="M 40 0 L 0 0 0 40" fill="none" stroke="white" stroke-width="1"/>
                    </pattern>
                </defs>
                <rect width="100%" height="100%" fill="url(#grid)" />
            </svg>
        </div>
        
        <!-- Gradient Overlay -->
        <div class="absolute inset-0 bg-gradient-to-br from-black/20 to-transparent"></div>

        <!-- Header Brand -->
        <div class="relative z-10">
            <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white/10 backdrop-blur-sm border border-white/20 rounded-lg flex items-center justify-center text-white">
                    <Package class="w-6 h-6" />
                </div>
                <span class="text-2xl font-bold tracking-tight">K-Mate</span>
            </div>
        </div>

        <!-- Central Visual Content: SLIDER AREA -->
        <div class="relative z-10 flex-1 flex flex-col justify-center items-center w-full">
            
            <div class="relative w-full max-w-lg h-[400px] flex items-center justify-center">
                <transition-group name="slide">
                    <div v-if="activeSlide === 0" key="slide1" class="absolute inset-0 flex flex-col items-center justify-center">
                        <!-- SLIDE 1: Ahli Bangunan -->
                        <div class="relative w-80 h-80 flex items-center justify-center mb-8">
                            <div class="absolute inset-0 hero-circle rounded-full filter blur-2xl"></div>
                            <div class="relative z-20 animate-float">
                                <div class="bg-white p-6 rounded-3xl shadow-2xl transform rotate-3 border-4 border-blue-50">
                                    <HardHat class="w-32 h-32 text-blue-700 stroke-[1.5]" />
                                </div>
                            </div>
                            <div class="absolute top-10 -left-4 z-10 animate-float-delayed">
                                <div class="bg-blue-600 p-4 rounded-2xl shadow-xl border border-blue-400/30">
                                    <Hammer class="w-10 h-10 text-white" />
                                </div>
                            </div>
                            <div class="absolute bottom-12 -right-6 z-30 animate-float">
                                <div class="bg-orange-500 p-4 rounded-2xl shadow-xl border border-orange-400/30 flex items-center gap-3">
                                    <Ruler class="w-6 h-6 text-white" />
                                </div>
                            </div>
                        </div>
                        <div class="text-center">
                            <h3 class="text-2xl font-bold mb-2">Mitra Profesional Konstruksi</h3>
                            <p class="text-blue-100 text-sm max-w-xs mx-auto">Dirancang khusus untuk kebutuhan kontraktor dan toko bangunan.</p>
                        </div>
                    </div>

                    <div v-if="activeSlide === 1" key="slide2" class="absolute inset-0 flex flex-col items-center justify-center">
                        <!-- SLIDE 2: Manajemen Stok -->
                        <div class="relative w-80 h-80 flex items-center justify-center mb-8">
                            <div class="absolute inset-0 hero-circle rounded-full filter blur-2xl"></div>
                            <div class="relative z-20 animate-float">
                                <div class="bg-white p-6 rounded-3xl shadow-2xl transform -rotate-3 border-4 border-blue-50">
                                    <ClipboardList class="w-32 h-32 text-blue-700 stroke-[1.5]" />
                                </div>
                            </div>
                            <div class="absolute top-4 right-0 z-10 animate-float-delayed">
                                <div class="bg-green-500 p-4 rounded-2xl shadow-xl border border-green-400/30">
                                    <PackageCheck class="w-8 h-8 text-white" />
                                </div>
                            </div>
                            <div class="absolute bottom-10 -left-6 z-30 animate-float">
                                <div class="bg-blue-500 p-4 rounded-2xl shadow-xl border border-blue-400/30">
                                    <BarChart3 class="w-8 h-8 text-white" />
                                </div>
                            </div>
                        </div>
                        <div class="text-center">
                            <h3 class="text-2xl font-bold mb-2">Kontrol Stok Terintegrasi</h3>
                            <p class="text-blue-100 text-sm max-w-xs mx-auto">Pantau inventaris material masuk dan keluar secara real-time.</p>
                        </div>
                    </div>

                    <div v-if="activeSlide === 2" key="slide3" class="absolute inset-0 flex flex-col items-center justify-center">
                        <!-- SLIDE 3: Logistik & Pengiriman -->
                        <div class="relative w-80 h-80 flex items-center justify-center mb-8">
                            <div class="absolute inset-0 hero-circle rounded-full filter blur-2xl"></div>
                            <div class="relative z-20 animate-float">
                                <div class="bg-white p-6 rounded-3xl shadow-2xl transform rotate-2 border-4 border-blue-50">
                                    <Truck class="w-32 h-32 text-blue-700 stroke-[1.5]" />
                                </div>
                            </div>
                            <div class="absolute top-12 -left-8 z-10 animate-float-delayed">
                                <div class="bg-yellow-500 p-4 rounded-2xl shadow-xl border border-yellow-400/30">
                                    <MapPin class="w-8 h-8 text-white" />
                                </div>
                            </div>
                        </div>
                        <div class="text-center">
                            <h3 class="text-2xl font-bold mb-2">Distribusi Material Cepat</h3>
                            <p class="text-blue-100 text-sm max-w-xs mx-auto">Atur jadwal pengiriman material ke lokasi proyek dengan efisien.</p>
                        </div>
                    </div>
                </transition-group>
            </div>

            <!-- Slider Indicators -->
            <div class="flex gap-3 mt-8 z-20">
                <button 
                    v-for="(slide, index) in slides" 
                    :key="index"
                    @click="setSlide(index)"
                    class="h-2 rounded-full transition-all duration-300"
                    :class="activeSlide === index ? 'w-8 bg-white' : 'w-2 bg-white/40 hover:bg-white/60'"
                ></button>
            </div>
        </div>

        <!-- Footer Quote -->
        <div class="relative z-10 hidden lg:flex absolute bottom-16 left-16 right-auto justify-between items-end border-t border-white/10 pt-6 max-w-xs">
            <div>
                <p class="text-sm font-medium text-blue-100">K-Mate System v2.0</p>
                <p class="text-xs text-blue-300 mt-1">Solusi Manajemen Toko Bangunan</p>
            </div>
        </div>
    </div>

    <!-- BAGIAN KANAN: Form Login -->
    <div class="w-full lg:w-1/2 flex items-center justify-center p-8 bg-white">
        <div class="w-full max-w-md space-y-8">
            
            <!-- Mobile Logo (Visible only on mobile) -->
            <div class="lg:hidden flex justify-center mb-8">
                <div class="flex items-center gap-2 text-blue-800">
                    <div class="w-10 h-10 bg-blue-800 text-white rounded-lg flex items-center justify-center">
                        <Package class="w-6 h-6" />
                    </div>
                    <span class="text-2xl font-bold">K-Mate</span>
                </div>
            </div>

            <!-- Greeting -->
            <div class="text-center lg:text-left">
                <h2 class="text-3xl font-bold text-gray-900 tracking-tight">Selamat Datang</h2>
                <p class="mt-2 text-gray-500">Silakan masuk untuk mengelola toko Anda.</p>
            </div>

            <!-- Form -->
            <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
                
                <!-- Alert Error (Conditional) -->
                <transition name="fade">
                  <div v-if="error" class="bg-red-50 border-l-4 border-red-500 p-4 rounded-md flex items-start gap-3">
                      <AlertCircle class="w-5 h-5 text-red-500 mt-0.5" />
                      <div>
                          <p class="text-sm text-red-700 font-medium">Gagal Masuk</p>
                          <p class="text-xs text-red-600">{{ error }}</p>
                      </div>
                  </div>
                </transition>

                <div class="space-y-5">
                    <!-- Username Field -->
                    <div>
                        <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400 z-10">
                                <User class="w-5 h-5" />
                            </div>
                            <BaseInput 
                                v-model="username"
                                id="username" 
                                name="username" 
                                type="text" 
                                required 
                                class="pl-10" 
                                placeholder="owner_toko"
                            />
                        </div>
                    </div>

                    <!-- Password Field -->
                    <div>
                        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400 z-10">
                                <Lock class="w-5 h-5" />
                            </div>
                            <BaseInput 
                                v-model="password"
                                :type="showPassword ? 'text' : 'password'" 
                                id="password" 
                                name="password" 
                                required 
                                class="pl-10 pr-10" 
                                placeholder="••••••••"
                            />
                            <button 
                                type="button" 
                                @click="togglePassword"
                                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 focus:outline-none"
                            >
                                <component :is="showPassword ? EyeOff : Eye" class="w-5 h-5" />
                            </button>
                        </div>
                        <div class="flex justify-end mt-2">
                            <a href="#" class="text-xs font-medium text-blue-600 hover:text-blue-500">Lupa password?</a>
                        </div>
                    </div>
                </div>

                <!-- Submit Button -->
                <div>
                    <BaseButton 
                        type="submit" 
                        :disabled="isLoading"
                        :loading="isLoading"
                        class="w-full bg-blue-800 hover:bg-blue-900 h-12 text-base shadow-lg hover:shadow-xl"
                    >
                        {{ isLoading ? 'Memproses...' : 'Masuk ke Dashboard' }}
                    </BaseButton>
                </div>

                <!-- Footer Help -->
                <div class="mt-6 text-center text-xs text-gray-500">
                    <p>Mengalami kendala? <a href="#" class="font-medium text-blue-600 hover:text-blue-500">Hubungi IT Support</a></p>
                    <p class="mt-2 text-gray-400">K-Mate v2.0 &copy; 2026</p>
                </div>
            </form>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import {
  Package,
  HardHat,
  Hammer,
  Ruler,
  ClipboardList,
  PackageCheck,
  BarChart3,
  Truck,
  MapPin,
  AlertCircle,
  User,
  Lock,
  Eye,
  EyeOff
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const error = ref('')

// Slider Logic
const activeSlide = ref(0)
const slides = [
    { id: 0, bgColor: 'bg-blue-800' },
    { id: 1, bgColor: 'bg-blue-900' },
    { id: 2, bgColor: 'bg-indigo-900' }
]
let slideInterval

const currentSlideData = computed(() => slides[activeSlide.value])

const nextSlide = () => {
    activeSlide.value = (activeSlide.value + 1) % slides.length
}

const setSlide = (index) => {
    activeSlide.value = index
    clearInterval(slideInterval)
    slideInterval = setInterval(nextSlide, 5000)
}

const togglePassword = () => {
    showPassword.value = !showPassword.value
}

const handleLogin = async () => {
    // Reset state
    error.value = ''
    isLoading.value = true

    // Call store action
    const result = await authStore.login(username.value, password.value)

    if (result.success) {
        // Redirect logic based on role
        const role = result.user.role
        if (role === 'owner') {
            router.push('/dashboard')
        } else {
            router.push('/pos')
        }
    } else {
        error.value = result.message
    }
    
    isLoading.value = false
}

onMounted(() => {
    slideInterval = setInterval(nextSlide, 5000)
})

onUnmounted(() => {
    if (slideInterval) clearInterval(slideInterval)
})
</script>

<style scoped>
/* Custom Animation for floating elements */
@keyframes float {
    0% { transform: translateY(0px); }
    50% { transform: translateY(-15px); }
    100% { transform: translateY(0px); }
}
@keyframes float-reverse {
    0% { transform: translateY(0px); }
    50% { transform: translateY(10px); }
    100% { transform: translateY(0px); }
}
.animate-float {
    animation: float 6s ease-in-out infinite;
}
.animate-float-delayed {
    animation: float-reverse 7s ease-in-out infinite 1s;
}

/* Professional Circle Background */
.hero-circle {
    background: radial-gradient(circle, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0) 70%);
}

/* Slider Transitions */
.slide-enter-active,
.slide-leave-active {
    transition: all 0.8s ease-in-out;
    position: absolute;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
}

.slide-enter-from {
    opacity: 0;
    transform: translateX(50px);
}

.slide-leave-to {
    opacity: 0;
    transform: translateX(-50px);
}

.slide-enter-to,
.slide-leave-from {
    opacity: 1;
    transform: translateX(0);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

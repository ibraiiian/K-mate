<template>
  <div class="fixed bottom-6 right-6 z-[50]">
    <!-- Popover Content -->
    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-2 opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100 scale-100"
      leave-to-class="translate-y-2 opacity-0 scale-95"
    >
      <div 
        v-if="isOpen" 
        class="absolute bottom-16 right-0 w-[380px] h-[500px] bg-white rounded-lg shadow-xl border border-slate-200 flex flex-col overflow-hidden"
      >
        <!-- Header: Minimalist Shadcn Style -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <span class="font-semibold text-sm text-slate-900">K-Mate Assistant</span>
            <span class="px-1.5 py-0.5 rounded-full bg-slate-100 text-[10px] font-medium text-slate-500">Beta</span>
          </div>
          <!-- Close Button (Ghost Variant) -->
          <button 
            @click="toggleChat" 
            class="h-6 w-6 rounded-md flex items-center justify-center text-slate-400 hover:text-slate-900 hover:bg-slate-100 transition"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Chat Area -->
        <div 
          ref="messagesContainer" 
          class="flex-1 overflow-y-auto p-4 space-y-4 bg-white"
        >
          <!-- Empty State / Welcome -->
          <div v-if="chatHistoryDisplay.length === 0" class="flex flex-col items-center justify-center h-full text-center p-6 text-slate-500">
            <div class="w-12 h-12 bg-slate-50 rounded-full flex items-center justify-center mb-3">
              <Sparkles class="w-5 h-5 text-slate-400" />
            </div>
            <h4 class="text-sm font-medium text-slate-900 mb-1">Ada yang bisa dibantu?</h4>
            <p class="text-xs text-slate-500 max-w-[200px]">
              Tanyakan tentang stok barang, omzet harian, atau insight toko Anda.
            </p>
          </div>

          <!-- Messages -->
          <div 
            v-for="(msg, index) in chatHistoryDisplay" 
            :key="index"
            class="flex flex-col gap-1 text-sm"
            :class="msg.role === 'user' ? 'items-end' : 'items-start'"
          >
            <!-- Label (Optional, maybe too cluttered) -->
            <span class="text-[10px] text-slate-400 px-1">
              {{ msg.role === 'user' ? 'Anda' : 'Assistant' }}
            </span>

            <div 
              class="px-3 py-2 rounded-md max-w-[85%]"
              :class="[
                msg.role === 'user' 
                  ? 'bg-slate-900 text-slate-50' 
                  : 'bg-slate-100 text-slate-800 border border-slate-200'
              ]"
            >
              <div v-if="msg.role === 'user'">{{ msg.content }}</div>
              <!-- Markdown Content for AI -->
              <div v-else class="prose prose-sm prose-slate max-w-none text-sm leading-relaxed" v-html="formatMarkdown(msg.content)"></div>
            </div>
          </div>

          <!-- Loading State (Skeleton) -->
          <div v-if="loading" class="flex flex-col items-start gap-1">
            <span class="text-[10px] text-slate-400 px-1">Assistant</span>
            <div class="bg-slate-100 rounded-md px-3 py-3 w-2/3 animate-pulse flex flex-col gap-2">
              <div class="h-2 bg-slate-200 rounded w-full"></div>
              <div class="h-2 bg-slate-200 rounded w-4/5"></div>
            </div>
          </div>
        </div>

        <!-- Footer / Input -->
        <div class="p-3 border-t border-slate-100 bg-white">
          <form @submit.prevent="sendMessage" class="flex items-center gap-2">
            <input 
              v-model="inputMessage" 
              type="text" 
              placeholder="Ketik sesuatu..." 
              class="flex-1 px-3 py-2 text-sm bg-transparent border border-slate-200 rounded-md placeholder:text-slate-400 focus:outline-none focus:ring-1 focus:ring-slate-900 focus:border-slate-900 transition"
              :disabled="loading"
            >
            <button 
              type="submit" 
              class="h-9 px-3 bg-slate-900 text-white rounded-md hover:bg-slate-800 disabled:opacity-50 disabled:cursor-not-allowed transition flex items-center justify-center"
              :disabled="!inputMessage.trim() || loading"
            >
              <ArrowUp class="w-4 h-4" />
            </button>
          </form>
          <div class="text-[10px] text-center text-slate-400 mt-2">
            AI dapat membuat kesalahan. Cek kembali informasi penting.
          </div>
        </div>
      </div>
    </transition>

    <!-- Trigger Button (Minimalist) -->
    <button 
      @click="toggleChat"
      class="h-12 w-12 bg-slate-900 hover:bg-slate-800 text-white rounded-full shadow-lg flex items-center justify-center transition-all hover:scale-105 active:scale-95 z-50"
    >
      <transition 
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 rotate-90"
        enter-to-class="opacity-100 rotate-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 rotate-0"
        leave-to-class="opacity-0 rotate-90"
        mode="out-in"
      >
        <X v-if="isOpen" class="w-5 h-5" />
        <MessageSquare v-else class="w-5 h-5" />
      </transition>
    </button>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { MessageSquare, X, ArrowUp, Sparkles } from 'lucide-vue-next'
import api from '@/api/axiosInstance'

const isOpen = ref(false)
const inputMessage = ref('')
const loading = ref(false)

// Chat History for backend (Groq/Gemini format)
const backendHistory = ref([]) 
// Chat History for display
const chatHistoryDisplay = ref([])

const messagesContainer = ref(null)

const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    scrollToBottom()
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatMarkdown = (text) => {
  if (!text) return ''
  // Basic Shadcn-friendly formatting
  return text
    .replace(/\*\*(.*?)\*\*/g, '<span class="font-semibold text-slate-900 text-content">$1</span>')
    .replace(/\*(.*?)\*/g, '<span class="italic text-content">$1</span>')
    .replace(/- /g, '• ') // Bullet points
    .replace(/\n/g, '<br>')
}

const sendMessage = async () => {
  if (!inputMessage.value.trim() || loading.value) return

  const userMsg = inputMessage.value
  inputMessage.value = ''

  // Add to display
  chatHistoryDisplay.value.push({ role: 'user', content: userMsg })
  
  loading.value = true
  scrollToBottom()

  try {
    const { data } = await api.post('/chat', {
      message: userMsg,
      history: backendHistory.value
    })

    if (data.status === 'success') {
      const response = data.data
      backendHistory.value = response.history
      chatHistoryDisplay.value.push({
        role: 'model',
        content: response.response
      })
    }
  } catch (error) {
    console.error('Chat error:', error)
    chatHistoryDisplay.value.push({
      role: 'model',
      content: 'Maaf, ada gangguan teknis. Coba lagi nanti.'
    })
  } finally {
    loading.value = false
    scrollToBottom()
  }
}
</script>

<style scoped>
/* Scoped overrides if Tailwind prose plugin is not enough */
:deep(.text-content) {
  display: inline;
}
</style>

<template>
  <div class="dashboard-card p-6 flex items-start gap-4">
    <div 
      class="p-3 rounded-lg flex items-center justify-center shrink-0"
      :class="colorClasses[color] || 'bg-gray-100 text-gray-600'"
    >
      <slot name="icon">
        <ChartBar class="w-5 h-5" />
      </slot>
    </div>
    
    <div class="flex-1 min-w-0">
      <p class="text-slate-500 text-xs font-bold uppercase tracking-wide mb-1">{{ title }}</p>
      <div v-if="loading" class="h-8 w-24 bg-gray-200 animate-pulse rounded"></div>
      <h3 v-else class="text-2xl font-bold text-slate-800 truncate">{{ value }}</h3>
      
      <div v-if="trend" class="mt-2 flex items-center">
        <span 
          class="text-xs font-bold px-2 py-0.5 rounded flex items-center border"
          :class="trendIsUp ? 'bg-green-50 text-green-600 border-green-100' : 'bg-red-50 text-red-600 border-red-100'"
        >
          <TrendingUp v-if="trendIsUp" class="w-3 h-3 mr-1" />
          <TrendingDown v-else class="w-3 h-3 mr-1" />
          {{ trend }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: String,
  value: String,
  color: {
    type: String,
    default: 'blue' // blue, red, green, orange, purple
  },
  trend: String,
  trendIsUp: {
    type: Boolean,
    default: true
  },
  loading: Boolean
})

import { TrendingUp, TrendingDown, ChartBar } from 'lucide-vue-next'

const colorClasses = {
  blue: 'bg-blue-50 text-blue-600 border border-blue-100',
  red: 'bg-red-50 text-red-600 border border-red-100',
  green: 'bg-green-50 text-green-600 border border-green-100',
  orange: 'bg-orange-50 text-orange-600 border border-orange-100',
  purple: 'bg-purple-50 text-purple-600 border border-purple-100',
  dark: 'bg-slate-800 text-white border border-slate-700'
}
</script>

<style scoped>
.dashboard-card {
    background: white;
    border: 1px solid #E2E8F0;
    border-radius: 0.75rem; /* rounded-xl */
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.dashboard-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0,0,0,0.04);
}
</style>

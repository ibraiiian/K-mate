<template>
  <button
    :class="[
      'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
      'h-10 px-4 py-2', // Standard size
      variantClasses[variant]
    ]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <div v-if="loading" class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"></div>
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'destructive', 'outline', 'secondary', 'ghost', 'link'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const variantClasses = {
  default: 'bg-slate-900 text-slate-50 hover:bg-slate-900/90',
  destructive: 'bg-red-500 text-slate-50 hover:bg-red-500/90',
  outline: 'border border-slate-200 bg-white hover:bg-slate-100 hover:text-slate-900',
  secondary: 'bg-slate-100 text-slate-900 hover:bg-slate-100/80',
  ghost: 'hover:bg-slate-100 hover:text-slate-900',
  link: 'text-slate-900 underline-offset-4 hover:underline'
}
</script>

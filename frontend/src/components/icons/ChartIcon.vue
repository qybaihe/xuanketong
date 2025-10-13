<template>
  <svg 
    :width="size" 
    :height="size" 
    viewBox="0 0 24 24" 
    fill="none" 
    xmlns="http://www.w3.org/2000/svg"
    :class="className"
  >
    <path 
      d="M3 3V21H21" 
      :stroke="color" 
      stroke-width="2" 
      stroke-linecap="round" 
      stroke-linejoin="round"
    />
    <path 
      d="M9 9L12 6L16 10L20 6" 
      :stroke="color" 
      stroke-width="2" 
      stroke-linecap="round" 
      stroke-linejoin="round"
    />
    <circle 
      cx="9" 
      cy="9" 
      r="1" 
      :fill="color"
    />
    <circle 
      cx="12" 
      cy="6" 
      r="1" 
      :fill="color"
    />
    <circle 
      cx="16" 
      cy="10" 
      r="1" 
      :fill="color"
    />
    <circle 
      cx="20" 
      cy="6" 
      r="1" 
      :fill="color"
    />
  </svg>
</template>

<script setup lang="ts">
interface Props {
  size?: number | string
  color?: string
  className?: string
}

withDefaults(defineProps<Props>(), {
  size: 24,
  color: 'currentColor',
  className: ''
})
</script>

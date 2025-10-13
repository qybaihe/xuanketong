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
      d="M20.59 13.41L13.42 20.58C13.2343 20.7663 13.0137 20.9138 12.7709 21.0141C12.5281 21.1145 12.2678 21.1662 12.005 21.1662C11.7422 21.1662 11.4819 21.1145 11.2391 21.0141C10.9963 20.9138 10.7757 20.7663 10.59 20.58L3.42 13.41C3.23425 13.2243 3.0868 13.0037 2.98649 12.7609C2.88618 12.5181 2.83447 12.2578 2.83447 11.995C2.83447 11.7322 2.88618 11.4719 2.98649 11.2291C3.0868 10.9863 3.23425 10.7657 3.42 10.58L10.59 3.41C10.7757 3.22425 10.9963 3.0768 11.2391 2.97649C11.4819 2.87618 11.7422 2.82447 12.005 2.82447C12.2678 2.82447 12.5281 2.87618 12.7709 2.97649C13.0137 3.0768 13.2343 3.22425 13.42 3.41L20.59 10.58C20.7757 10.7657 20.9232 10.9863 21.0235 11.2291C21.1238 11.4719 21.1755 11.7322 21.1755 11.995C21.1755 12.2578 21.1238 12.5181 21.0235 12.7609C20.9232 13.0037 20.7757 13.2243 20.59 13.41Z" 
      :stroke="color" 
      stroke-width="2" 
      stroke-linecap="round" 
      stroke-linejoin="round"
    />
    <path 
      d="M7 7H7.01" 
      :stroke="color" 
      stroke-width="2" 
      stroke-linecap="round" 
      stroke-linejoin="round"
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

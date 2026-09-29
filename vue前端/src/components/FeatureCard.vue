<template>
  <div 
    class="feature-card relative overflow-hidden rounded-xl p-6 border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shadow-md hover:shadow-xl transition-all duration-300"
    :class="[`feature-card-${color}`]"
    v-motion
    :initial="{ opacity: 0, y: 100 }"
    :enter="{ opacity: 1, y: 0, transition: { delay: delay, duration: 800 } }"
  >
    <div class="absolute top-0 right-0 w-20 h-20 -mr-10 -mt-10 rounded-full bg-gradient-to-br opacity-20"
      :class="[`from-${color}-400 to-${color}-600`]"></div>
    
    <div class="relative z-10">
      <div class="flex items-center mb-4">
        <div class="flex items-center justify-center h-12 w-12 rounded-md text-white"
          :class="[`bg-${color}-500 dark:bg-${color}-600`]">
          <slot name="icon"></slot>
        </div>
        <h3 class="ml-4 text-xl font-bold text-gray-900 dark:text-white">{{ title }}</h3>
      </div>
      
      <p class="text-gray-600 dark:text-gray-300">{{ description }}</p>
      
      <div class="mt-4">
        <slot></slot>
      </div>
    </div>
    
    <div class="absolute bottom-0 right-0 w-16 h-16 opacity-10">
      <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
        <path fill="currentColor" d="M45.3,-51.2C58.3,-40.9,68.5,-25.4,71.8,-8.1C75.1,9.2,71.5,28.3,60.7,41.5C49.9,54.7,32,62,12.4,68.2C-7.2,74.4,-28.5,79.4,-43.7,71.8C-58.9,64.1,-68,43.7,-72.1,23.1C-76.2,2.5,-75.3,-18.3,-66.6,-34.1C-57.9,-49.9,-41.4,-60.7,-24.8,-68.8C-8.2,-76.9,8.5,-82.3,23.3,-77.1C38.1,-71.9,51,-61.1,45.3,-51.2Z" transform="translate(100 100)" />
      </svg>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  color: {
    type: String,
    default: 'blue'
  },
  delay: {
    type: Number,
    default: 0
  }
})
</script>

<style scoped>
.feature-card {
  backdrop-filter: blur(5px);
}

.feature-card-blue {
  border-top: 3px solid #3b82f6;
}

.feature-card-indigo {
  border-top: 3px solid #6366f1;
}

.feature-card-purple {
  border-top: 3px solid #8b5cf6;
}

.feature-card-green {
  border-top: 3px solid #10b981;
}

.feature-card-red {
  border-top: 3px solid #ef4444;
}

.feature-card-orange {
  border-top: 3px solid #f97316;
}
</style> 
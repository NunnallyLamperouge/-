<template>
  <div class="wave-container relative overflow-hidden" :style="{ height: height }">
    <div class="absolute inset-0" :class="bgClass"></div>
    
    <div class="wave-wrapper absolute bottom-0 left-0 w-full overflow-hidden">
      <svg 
        class="wave" 
        :class="{ 'wave-animated': animated }"
        viewBox="0 0 1440 120" 
        version="1.1" 
        xmlns="http://www.w3.org/2000/svg"
      >
        <path 
          :fill="waveColor" 
          fill-opacity="0.8"
          d="M0,0 C240,40 480,80 720,80 C960,80 1200,40 1440,0 L1440,120 L0,120 Z"
        ></path>
      </svg>
      <svg 
        class="wave wave-2" 
        :class="{ 'wave-animated': animated }"
        viewBox="0 0 1440 120" 
        version="1.1" 
        xmlns="http://www.w3.org/2000/svg"
      >
        <path 
          :fill="waveColor" 
          fill-opacity="0.3"
          d="M0,0 C240,60 480,100 720,80 C960,60 1200,20 1440,40 L1440,120 L0,120 Z"
        ></path>
      </svg>
    </div>
    
    <div class="relative z-10 h-full">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  height: {
    type: String,
    default: '400px'
  },
  bgClass: {
    type: String,
    default: 'bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-700 dark:to-indigo-700'
  },
  waveColor: {
    type: String,
    default: '#ffffff'
  },
  animated: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.wave-container {
  position: relative;
}

.wave-wrapper {
  height: 120px;
  transform: translateY(1px); /* 防止波浪和下一个部分之间出现缝隙 */
}

.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 120px;
  transform: translateX(0);
}

.wave-2 {
  animation-delay: -0.5s;
  opacity: 0.5;
}

.wave-animated {
  animation: wave 15s linear infinite;
}

@keyframes wave {
  0% {
    transform: translateX(0);
  }
  50% {
    transform: translateX(-50%);
  }
  100% {
    transform: translateX(0);
  }
}
</style> 
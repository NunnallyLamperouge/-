<template>
  <div class="virtual-defense-simulator relative rounded-lg overflow-hidden shadow-xl" :style="{ height: height }">
    <!-- 模拟3D场景 -->
    <div class="absolute inset-0 bg-gradient-to-br from-gray-900 to-gray-800">
      <!-- 网格地面 -->
      <div class="absolute inset-0 perspective-1000">
        <div class="grid-floor"></div>
      </div>
      
      <!-- 虚拟人物 -->
      <div class="absolute bottom-0 left-1/2 transform -translate-x-1/2 flex space-x-16 items-end">
        <!-- 答辩委员 -->
        <div class="virtual-person committee" v-for="i in 3" :key="`committee-${i}`">
          <div class="head"></div>
          <div class="body"></div>
          <div class="name-tag">答辩委员 {{ i }}</div>
        </div>
        
        <!-- 学生 -->
        <div class="virtual-person student">
          <div class="head"></div>
          <div class="body"></div>
          <div class="name-tag">答辩学生</div>
        </div>
      </div>
      
      <!-- 对话气泡 -->
      <div class="speech-bubble" v-motion :initial="{ opacity: 0, scale: 0.8 }" :enter="{ opacity: 1, scale: 1, transition: { delay: 500, duration: 500 } }">
        <div class="bubble-header">
          <span class="bubble-title">答辩问题</span>
          <span class="bubble-badge">难度: 中等</span>
        </div>
        <div class="bubble-content">
          请解释一下您论文中提到的"动态样式映射算法"是如何实现不同格式标准之间的自动转换的？
        </div>
      </div>
      
      <!-- 控制面板 -->
      <div class="control-panel">
        <div class="panel-header">虚拟答辩控制台</div>
        <div class="panel-content">
          <div class="control-item">
            <span class="control-label">难度级别</span>
            <div class="control-buttons">
              <button class="control-button">初级</button>
              <button class="control-button active">中级</button>
              <button class="control-button">高级</button>
            </div>
          </div>
          <div class="control-item">
            <span class="control-label">委员数量</span>
            <div class="control-slider">
              <div class="slider-track">
                <div class="slider-fill"></div>
                <div class="slider-thumb"></div>
              </div>
              <div class="slider-value">3</div>
            </div>
          </div>
          <div class="control-item">
            <span class="control-label">情绪分析</span>
            <div class="emotion-meter">
              <div class="emotion-fill" style="width: 65%"></div>
            </div>
            <span class="emotion-value">65% 自信</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  height: {
    type: String,
    default: '400px'
  }
})
</script>

<style scoped>
.virtual-defense-simulator {
  position: relative;
}

/* 网格地面 */
.grid-floor {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 60%;
  background-image: 
    linear-gradient(to right, rgba(50, 50, 50, 0.3) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(50, 50, 50, 0.3) 1px, transparent 1px);
  background-size: 40px 40px;
  transform: rotateX(60deg);
  transform-origin: bottom;
}

/* 虚拟人物 */
.virtual-person {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  margin-bottom: 20px;
}

.head {
  width: 20px;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 50%;
}

.body {
  width: 30px;
  height: 60px;
  background-color: #3b82f6;
  border-radius: 5px;
  margin-top: 5px;
}

.committee .body {
  background-color: #8b5cf6;
}

.name-tag {
  position: absolute;
  bottom: -20px;
  white-space: nowrap;
  font-size: 10px;
  color: #f0f0f0;
  background-color: rgba(0, 0, 0, 0.5);
  padding: 2px 5px;
  border-radius: 3px;
}

/* 对话气泡 */
.speech-bubble {
  position: absolute;
  top: 20%;
  left: 50%;
  transform: translateX(-50%);
  width: 80%;
  max-width: 400px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 10px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.bubble-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #3b82f6;
  color: white;
}

.bubble-badge {
  font-size: 12px;
  background-color: rgba(0, 0, 0, 0.2);
  padding: 2px 8px;
  border-radius: 10px;
}

.bubble-content {
  padding: 15px;
  color: #333;
}

/* 控制面板 */
.control-panel {
  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 200px;
  background-color: rgba(30, 30, 30, 0.8);
  border-radius: 10px;
  overflow: hidden;
  backdrop-filter: blur(5px);
}

.panel-header {
  padding: 8px 12px;
  background-color: #3b82f6;
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.panel-content {
  padding: 10px;
}

.control-item {
  margin-bottom: 10px;
}

.control-label {
  display: block;
  font-size: 11px;
  color: #ccc;
  margin-bottom: 5px;
}

.control-buttons {
  display: flex;
  gap: 5px;
}

.control-button {
  flex: 1;
  font-size: 10px;
  padding: 3px 0;
  background-color: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 3px;
  color: #ccc;
  cursor: pointer;
}

.control-button.active {
  background-color: #3b82f6;
  color: white;
}

.control-slider {
  display: flex;
  align-items: center;
  gap: 10px;
}

.slider-track {
  flex: 1;
  height: 4px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  position: relative;
}

.slider-fill {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 60%;
  background-color: #3b82f6;
  border-radius: 2px;
}

.slider-thumb {
  position: absolute;
  left: 60%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 12px;
  height: 12px;
  background-color: white;
  border-radius: 50%;
  cursor: pointer;
}

.slider-value {
  font-size: 12px;
  color: white;
  width: 20px;
  text-align: center;
}

.emotion-meter {
  height: 6px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  margin-bottom: 5px;
  overflow: hidden;
}

.emotion-fill {
  height: 100%;
  background: linear-gradient(to right, #10b981, #3b82f6);
  border-radius: 3px;
}

.emotion-value {
  font-size: 10px;
  color: #ccc;
}
</style> 
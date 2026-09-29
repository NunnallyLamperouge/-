<template>
  <div class="semantic-enhancer-preview relative rounded-lg overflow-hidden shadow-xl" :style="{ height: height }">
    <div class="absolute inset-0 bg-white dark:bg-gray-900">
      <!-- 编辑器界面 -->
      <div class="h-full flex flex-col">
        <!-- 工具栏 -->
        <div class="bg-gray-100 dark:bg-gray-800 px-4 py-2 flex items-center justify-between border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center space-x-4">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300">语义增强编辑器</div>
            <div class="flex items-center space-x-2">
              <button class="tool-button">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                </svg>
              </button>
              <button class="tool-button">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2h-1V9a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </button>
              <button class="tool-button">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M17.414 2.586a2 2 0 00-2.828 0L7 10.172V13h2.828l7.586-7.586a2 2 0 000-2.828z" />
                  <path fill-rule="evenodd" d="M2 6a2 2 0 012-2h4a1 1 0 010 2H4v10h10v-4a1 1 0 112 0v4a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <button class="enhance-button">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clip-rule="evenodd" />
              </svg>
              增强
            </button>
          </div>
        </div>
        
        <!-- 编辑区域 -->
        <div class="flex-1 flex">
          <!-- 文本编辑区 -->
          <div class="w-2/3 p-4 overflow-auto editor-area">
            <p class="paragraph">
              <span class="text-normal">在本研究中，我们</span>
              <span class="text-error" data-error="语法错误：主谓不一致">提出了一种新方法，并且验证了它的有效性</span>
              <span class="text-normal">。该方法基于深度学习技术，</span>
              <span class="text-warning" data-suggestion="建议使用更专业的表达：'通过卷积神经网络提取特征'">使用神经网络来处理数据</span>
              <span class="text-normal">，从而实现了更高的准确率。</span>
            </p>
            
            <p class="paragraph">
              <span class="text-normal">实验结果表明，与传统方法相比，</span>
              <span class="text-improvement" data-improvement="可以更具体：'本文提出的方法在准确率上提高了15.3%，在处理速度上提升了2.7倍'">我们的方法有明显的优势</span>
              <span class="text-normal">。特别是在处理大规模数据集时，性能提升更为显著。</span>
            </p>
            
            <p class="paragraph">
              <span class="text-normal">此外，</span>
              <span class="text-error" data-error="逻辑错误：前后矛盾">虽然我们的方法计算复杂度较低，但是运行速度却比基准方法慢</span>
              <span class="text-normal">。这主要是因为我们采用了更复杂的模型结构。</span>
            </p>
            
            <p class="paragraph">
              <span class="text-normal">在未来的工作中，我们计划</span>
              <span class="text-warning" data-suggestion="过于笼统，建议：'优化模型架构，引入注意力机制，并探索模型压缩技术'">进一步改进我们的方法</span>
              <span class="text-normal">，以提高其在实际应用中的表现。</span>
            </p>
          </div>
          
          <!-- 建议面板 -->
          <div class="w-1/3 bg-gray-50 dark:bg-gray-800 border-l border-gray-200 dark:border-gray-700 overflow-auto">
            <div class="p-4">
              <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">增强建议</h3>
              
              <div class="suggestion-item error" v-motion :initial="{ opacity: 0, x: 20 }" :enter="{ opacity: 1, x: 0, transition: { delay: 200 } }">
                <div class="suggestion-header">
                  <div class="suggestion-badge error">错误</div>
                  <div class="suggestion-text">语法错误：主谓不一致</div>
                </div>
                <div class="suggestion-content">
                  <p class="text-xs text-gray-600 dark:text-gray-400 mb-2">原文：提出了一种新方法，并且验证了它的有效性</p>
                  <p class="text-xs text-gray-800 dark:text-gray-200">建议：提出了一种新方法，并且验证了其有效性</p>
                </div>
                <div class="suggestion-actions">
                  <button class="action-button">应用</button>
                  <button class="action-button secondary">忽略</button>
                </div>
              </div>
              
              <div class="suggestion-item warning" v-motion :initial="{ opacity: 0, x: 20 }" :enter="{ opacity: 1, x: 0, transition: { delay: 400 } }">
                <div class="suggestion-header">
                  <div class="suggestion-badge warning">优化</div>
                  <div class="suggestion-text">表达不够专业</div>
                </div>
                <div class="suggestion-content">
                  <p class="text-xs text-gray-600 dark:text-gray-400 mb-2">原文：使用神经网络来处理数据</p>
                  <p class="text-xs text-gray-800 dark:text-gray-200">建议：通过卷积神经网络提取特征</p>
                </div>
                <div class="suggestion-actions">
                  <button class="action-button">应用</button>
                  <button class="action-button secondary">忽略</button>
                </div>
              </div>
              
              <div class="suggestion-item error" v-motion :initial="{ opacity: 0, x: 20 }" :enter="{ opacity: 1, x: 0, transition: { delay: 600 } }">
                <div class="suggestion-header">
                  <div class="suggestion-badge error">错误</div>
                  <div class="suggestion-text">逻辑错误：前后矛盾</div>
                </div>
                <div class="suggestion-content">
                  <p class="text-xs text-gray-600 dark:text-gray-400 mb-2">原文：虽然我们的方法计算复杂度较低，但是运行速度却比基准方法慢</p>
                  <p class="text-xs text-gray-800 dark:text-gray-200">建议：虽然我们的方法理论计算复杂度较高，但通过优化实现，实际运行速度比基准方法快</p>
                </div>
                <div class="suggestion-actions">
                  <button class="action-button">应用</button>
                  <button class="action-button secondary">忽略</button>
                </div>
              </div>
              
              <div class="suggestion-item improvement" v-motion :initial="{ opacity: 0, x: 20 }" :enter="{ opacity: 1, x: 0, transition: { delay: 800 } }">
                <div class="suggestion-header">
                  <div class="suggestion-badge improvement">改进</div>
                  <div class="suggestion-text">表达过于笼统</div>
                </div>
                <div class="suggestion-content">
                  <p class="text-xs text-gray-600 dark:text-gray-400 mb-2">原文：我们的方法有明显的优势</p>
                  <p class="text-xs text-gray-800 dark:text-gray-200">建议：本文提出的方法在准确率上提高了15.3%，在处理速度上提升了2.7倍</p>
                </div>
                <div class="suggestion-actions">
                  <button class="action-button">应用</button>
                  <button class="action-button secondary">忽略</button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 状态栏 -->
        <div class="bg-gray-100 dark:bg-gray-800 px-4 py-2 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center text-xs text-gray-600 dark:text-gray-400">
          <div class="flex items-center space-x-4">
            <div>错误: <span class="text-red-500">2</span></div>
            <div>警告: <span class="text-yellow-500">1</span></div>
            <div>建议: <span class="text-blue-500">1</span></div>
          </div>
          <div>
            学术表达强度: <span class="text-green-500">78%</span>
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
.tool-button {
  @apply p-1 rounded text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700;
}

.enhance-button {
  @apply flex items-center px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600;
}

.editor-area {
  font-family: 'Times New Roman', serif;
  line-height: 1.8;
}

.paragraph {
  margin-bottom: 1.5em;
}

.text-normal {
  @apply text-gray-800 dark:text-gray-200;
}

.text-error {
  @apply text-red-600 dark:text-red-400 bg-red-100 dark:bg-red-900 dark:bg-opacity-30 px-0.5 rounded border-b-2 border-red-500 cursor-help;
}

.text-warning {
  @apply text-yellow-600 dark:text-yellow-400 bg-yellow-100 dark:bg-yellow-900 dark:bg-opacity-30 px-0.5 rounded border-b-2 border-yellow-500 cursor-help;
}

.text-improvement {
  @apply text-blue-600 dark:text-blue-400 bg-blue-100 dark:bg-blue-900 dark:bg-opacity-30 px-0.5 rounded border-b-2 border-blue-500 cursor-help;
}

.suggestion-item {
  @apply mb-4 rounded-lg overflow-hidden border;
}

.suggestion-item.error {
  @apply border-red-200 dark:border-red-900;
}

.suggestion-item.warning {
  @apply border-yellow-200 dark:border-yellow-900;
}

.suggestion-item.improvement {
  @apply border-blue-200 dark:border-blue-900;
}

.suggestion-header {
  @apply flex items-center p-2 bg-gray-100 dark:bg-gray-700;
}

.suggestion-badge {
  @apply text-xs px-2 py-0.5 rounded-full mr-2 font-medium;
}

.suggestion-badge.error {
  @apply bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200;
}

.suggestion-badge.warning {
  @apply bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200;
}

.suggestion-badge.improvement {
  @apply bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200;
}

.suggestion-text {
  @apply text-xs font-medium text-gray-700 dark:text-gray-300;
}

.suggestion-content {
  @apply p-2 bg-white dark:bg-gray-800;
}

.suggestion-actions {
  @apply flex justify-end p-2 bg-gray-50 dark:bg-gray-700;
}

.action-button {
  @apply text-xs px-2 py-1 rounded bg-blue-500 text-white hover:bg-blue-600 mr-2;
}

.action-button.secondary {
  @apply bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-500;
}
</style> 
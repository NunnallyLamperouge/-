<template>
  <div class="smart-formatting-preview relative rounded-lg overflow-hidden shadow-xl" :style="{ height: height }">
    <div class="absolute inset-0 bg-white dark:bg-gray-900">
      <!-- 分屏显示 -->
      <div class="flex h-full">
        <!-- 左侧：原始文档 -->
        <div class="w-1/2 h-full border-r border-gray-200 dark:border-gray-700 relative">
          <div class="absolute top-0 left-0 right-0 bg-gray-100 dark:bg-gray-800 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 flex justify-between items-center">
            <span>原始文档</span>
            <div class="flex items-center space-x-2">
              <span class="text-xs text-gray-500 dark:text-gray-400">Word 格式</span>
              <div class="w-3 h-3 rounded-full bg-red-500"></div>
            </div>
          </div>
          
          <div class="p-4 pt-12 h-full overflow-auto">
            <div class="doc-content original-doc">
              <div class="doc-title">基于深度学习的智能排版系统研究</div>
              <div class="doc-author">张三</div>
              <div class="doc-abstract-title">摘要</div>
              <div class="doc-abstract-content">
                本文提出了一种基于深度学习的智能排版系统，该系统能够自动识别文档结构并根据不同的排版规范进行格式转换。实验结果表明，该系统能够有效减少人工排版的时间成本，提高排版质量。
              </div>
              <div class="doc-keywords">关键词：深度学习；智能排版；格式转换</div>
              <div class="doc-section-title">1. 引言</div>
              <div class="doc-paragraph">
                随着学术研究的不断发展，各高校和期刊对论文格式的要求也越来越规范化。然而，不同机构之间的格式规范存在较大差异，这给研究人员带来了额外的工作负担。
              </div>
              <div class="doc-paragraph">
                传统的排版方法主要依靠人工调整，不仅耗时耗力，而且容易出错。因此，开发一种能够自动适应不同格式规范的智能排版系统具有重要意义。
              </div>
            </div>
          </div>
        </div>
        
        <!-- 右侧：转换后文档 -->
        <div class="w-1/2 h-full relative">
          <div class="absolute top-0 left-0 right-0 bg-blue-50 dark:bg-blue-900 px-4 py-2 text-sm font-medium text-blue-700 dark:text-blue-300 flex justify-between items-center">
            <span>转换后文档</span>
            <div class="flex items-center space-x-2">
              <span class="text-xs text-blue-500 dark:text-blue-400">IEEE 格式</span>
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
            </div>
          </div>
          
          <div class="p-4 pt-12 h-full overflow-auto">
            <div class="doc-content formatted-doc" v-motion :initial="{ opacity: 0 }" :enter="{ opacity: 1, transition: { delay: 300, duration: 800 } }">
              <div class="ieee-title">基于深度学习的智能排版系统研究</div>
              <div class="ieee-author">张三</div>
              <div class="ieee-affiliation">计算机科学与技术学院，某某大学，中国</div>
              <div class="ieee-abstract-title">摘要</div>
              <div class="ieee-abstract-content">
                本文提出了一种基于深度学习的智能排版系统，该系统能够自动识别文档结构并根据不同的排版规范进行格式转换。实验结果表明，该系统能够有效减少人工排版的时间成本，提高排版质量。
              </div>
              <div class="ieee-keywords"><span class="font-bold">关键词</span>—深度学习；智能排版；格式转换</div>
              <div class="ieee-section-title">I. 引言</div>
              <div class="ieee-paragraph">
                随着学术研究的不断发展，各高校和期刊对论文格式的要求也越来越规范化。然而，不同机构之间的格式规范存在较大差异，这给研究人员带来了额外的工作负担。
              </div>
              <div class="ieee-paragraph">
                传统的排版方法主要依靠人工调整，不仅耗时耗力，而且容易出错。因此，开发一种能够自动适应不同格式规范的智能排版系统具有重要意义。
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 转换过程动画 -->
      <div class="conversion-animation">
        <svg width="60" height="30" viewBox="0 0 60 30" xmlns="http://www.w3.org/2000/svg">
          <path class="arrow-path" d="M5,15 L55,15" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" />
          <path class="arrow-head" d="M45,5 L55,15 L45,25" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
          <circle class="moving-dot" cx="5" cy="15" r="4" fill="#3b82f6" />
        </svg>
      </div>
      
      <!-- 底部工具栏 -->
      <div class="absolute bottom-0 left-0 right-0 bg-gray-100 dark:bg-gray-800 px-4 py-2 flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <div class="template-selector flex items-center space-x-2">
            <span class="text-xs text-gray-600 dark:text-gray-400">目标格式:</span>
            <select class="text-xs bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded px-2 py-1 text-gray-700 dark:text-gray-300">
              <option>IEEE</option>
              <option>ACM</option>
              <option>北京大学</option>
              <option>清华大学</option>
              <option>自定义...</option>
            </select>
          </div>
          
          <button class="convert-button text-xs bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded flex items-center space-x-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
            </svg>
            <span>转换</span>
          </button>
        </div>
        
        <div class="conversion-stats text-xs text-gray-600 dark:text-gray-400">
          转换用时: <span class="font-medium text-blue-600 dark:text-blue-400">0.3秒</span> | 
          匹配度: <span class="font-medium text-green-600 dark:text-green-400">98%</span>
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
/* 文档样式 */
.doc-content {
  font-family: 'Times New Roman', serif;
  line-height: 1.6;
}

/* 原始文档样式 */
.original-doc .doc-title {
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 10px;
}

.original-doc .doc-author {
  font-size: 14px;
  text-align: center;
  margin-bottom: 20px;
}

.original-doc .doc-abstract-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 5px;
}

.original-doc .doc-abstract-content {
  font-size: 14px;
  margin-bottom: 10px;
  text-indent: 2em;
}

.original-doc .doc-keywords {
  font-size: 14px;
  margin-bottom: 20px;
}

.original-doc .doc-section-title {
  font-size: 16px;
  font-weight: bold;
  margin: 15px 0 10px;
}

.original-doc .doc-paragraph {
  font-size: 14px;
  margin-bottom: 10px;
  text-indent: 2em;
}

/* IEEE格式文档样式 */
.ieee-title {
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 10px;
}

.ieee-author {
  font-size: 14px;
  text-align: center;
  margin-bottom: 5px;
}

.ieee-affiliation {
  font-size: 12px;
  font-style: italic;
  text-align: center;
  margin-bottom: 20px;
}

.ieee-abstract-title {
  font-size: 14px;
  font-weight: bold;
  font-style: italic;
  margin-bottom: 5px;
}

.ieee-abstract-content {
  font-size: 12px;
  margin-bottom: 10px;
}

.ieee-keywords {
  font-size: 12px;
  font-style: italic;
  margin-bottom: 20px;
}

.ieee-section-title {
  font-size: 14px;
  font-weight: bold;
  margin: 15px 0 10px;
}

.ieee-paragraph {
  font-size: 12px;
  margin-bottom: 10px;
  text-indent: 0;
}

/* 转换动画 */
.conversion-animation {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 30px;
  z-index: 10;
}

.arrow-path {
  stroke-dasharray: 50;
  stroke-dashoffset: 0;
  animation: dash 2s linear infinite;
}

.moving-dot {
  animation: moveDot 2s linear infinite;
}

@keyframes dash {
  from {
    stroke-dashoffset: 50;
  }
  to {
    stroke-dashoffset: 0;
  }
}

@keyframes moveDot {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(50px);
  }
}
</style> 
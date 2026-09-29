<script setup>
import { LoadingOutlined } from "@ant-design/icons-vue";
import { h } from "vue";
const indicator = h(LoadingOutlined, {
  style: {
    fontSize: "36px",
  },
  spin: true,
});
</script>
<template>
  <div class="flex flex-col items-center justify-center">
    <a-spin class="text-green-100" :indicator="indicator" />
    <h1 class="text-xl dark:text-white mt-2">加载中...</h1>
  </div>
</template>

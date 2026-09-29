<script setup>
import { LoadingOutlined } from "@ant-design/icons-vue";
import { h } from "vue";
const indicator = h(LoadingOutlined, {
  style: {
    fontSize: "18px",
  },
  spin: true,
});
</script>
<template>
  <a-spin class="text-green-100" :indicator="indicator" />
</template>

<template>
  <router-view></router-view>
</template>

<script lang="ts">
export default {
  name: 'App'
}
</script>

<style>
/* Element Plus themed styles are handled globally */

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Arial, sans-serif;
  line-height: 1.6;
  color: #333;
}

.main-container {
  display: flex;
  min-height: 100vh;
}

.content {
  flex: 1;
  padding: 20px;
}

.card {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card h3 {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  color: #333;
}

.card h3 i {
  cursor: pointer;
  color: #4a90e2;
  transition: transform 0.3s;
}

.card h3 i:hover {
  transform: rotate(180deg);
}
</style>

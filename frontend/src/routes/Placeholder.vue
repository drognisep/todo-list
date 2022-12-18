<template>
<div class="container">
  <h1>This is a route placeholder</h1>
</div>
</template>

<script>
export default {
  name: "PlaceholderRoute"
}
</script>

<style scoped>
.container {
  position: relative;
  width: 100%;
  height: 100%;
}

.container h1 {
  line-height: calc(100vh - var(--toolbar-height));
  text-align: center;
  margin: auto 0;
}
</style>
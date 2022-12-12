<template>
  <div class="overlay" @click.prevent.stop>
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: "Overlay",
}
</script>

<style scoped>
.overlay {
  position: fixed;
  width: 100vw;
  height: 100vh;
  top: 0;
  left: 0;
  background-color: var(--overlay-color);
  z-index: var(--z-overlay);
}
</style>
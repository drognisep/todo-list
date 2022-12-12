<template>
  <overlay>
    <img alt="loading" class="loading" src="../assets/images/gear.svg"/>
  </overlay>
</template>

<script>
import Overlay from "./Overlay.vue";

export default {
  name: "Loading",
  components: {Overlay},
}
</script>

<style scoped>
.loading {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100px;
  height: 100px;
  margin: -50px 0 0 -50px;
  animation: rotation 1s infinite linear;
}

@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(359deg);
  }
}
</style>
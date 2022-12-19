<template>
  <div id="main-view">
    <router-view></router-view>
  </div>
  <Progress/>
  <Confirm/>
</template>

<script>
import {RouterView} from 'vue-router';
import Confirm from "./Confirm.vue";
import Progress from "./Progress.vue";

export default {
  name: "MainView",
  components: {RouterView, Confirm, Progress},
}
</script>

<style scoped>
#main-view {
  position: absolute;
  top: var(--toolbar-height);
  left: 0;
  right: 0;
  bottom: var(--footer-height);
  overflow: hidden;
}
</style>
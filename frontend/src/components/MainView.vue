<template>
  <div id="main-view">
    <router-view></router-view>
  </div>
</template>

<script>
import {RouterView} from 'vue-router';

export default {
  name: "MainView",
  components: [RouterView],
}
</script>

<style scoped>
#main-view {
  position: absolute;
  top: var(--toolbar-height);
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}
</style>
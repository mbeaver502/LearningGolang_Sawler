<template>
  <Header />
  <div>
    <router-view />
  </div>
  <Footer />
</template>

<script>
import PageHeader from "./components/PageHeader.vue";
import PageFooter from "./components/PageFooter.vue";
import { store } from "./components/store.js";

export default {
  name: "App",
  components: {
    Header: PageHeader,
    Footer: PageFooter,
  },
  data() {
    return {
      store,
    };
  },
};
</script>

<style>
</style>
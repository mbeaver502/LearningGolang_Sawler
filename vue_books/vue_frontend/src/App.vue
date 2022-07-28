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

export default {
  name: "App",
  components: {
    Header: PageHeader,
    Footer: PageFooter,
  },
};
</script>

<style>
</style>
<template>
  <Header />
  <Body />
  <Footer />
</template>

<script>
import PageHeader from "./components/PageHeader.vue";
import PageBody from "./components/PageBody.vue";
import PageFooter from "./components/PageFooter.vue";

export default {
  name: "App",
  components: {
    Header: PageHeader,
    Body: PageBody,
    Footer: PageFooter,
  },
};
</script>

<style>
</style>
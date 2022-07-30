<template>
  <div class="container">
    <div v-if="ready" class="row">
      <div class="col-md-2">
        <img
          class="img-fluid img-thumbnail"
          :src="`${imgPath}/covers/${book.slug}.jpg`"
          alt="cover"
        />
      </div>
      <div class="col-md-10">
        <h3 class="mt-3">
          {{ book.title }}
        </h3>
        <hr />
        <p>
          <strong>Author: </strong>{{ book.author.author_name }}<br />
          <strong>Published: </strong>{{ book.publication_year }}
        </p>
        <p>
          {{ book.description }}
        </p>
      </div>
    </div>
    <p v-else>Loading...</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      book: {},
      imgPath: process.env.VUE_APP_IMAGE_URL,
      ready: false,
    };
  },
  // ensure this keep-alive'd component gets the right data on page load
  activated() {
    fetch(process.env.VUE_APP_API_URL + "/books/" + this.$route.params.bookName)
      .then((response) => response.json())
      .then((response) => {
        if (response.error) {
          this.$emit("error", response.message);
        } else {
          this.book = response.data;
          this.ready = true;
        }
      });
  },
};
</script>
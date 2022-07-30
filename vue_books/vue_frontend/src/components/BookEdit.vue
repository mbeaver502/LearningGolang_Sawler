<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-3">Add / Edit Book</h1>
      </div>

      <hr />

      <form-tag
        @bookEditEvent="submitHandler"
        name="bookForm"
        event="bookEditEvent"
      >
        <div v-if="this.book.slug !== ''" class="mb-3">
          <img
            :src="`${this.imgPath}/covers/${this.book.slug}.jpg`"
            class="img-fluid img-thumbnail book-cover"
            alt="cover"
          />
        </div>

        <div class="mb-3">
          <label for="formFile" class="form-label">Cover Image</label>
          <input
            v-if="this.book.id === 0"
            ref="coverInput"
            class="form-control"
            type="file"
            id="formFile"
            required
            accept="image/jpeg"
            @change="loadCoverImage"
          />
          <input
            v-else
            ref="coverInput"
            class="form-control"
            type="file"
            id="formFile"
            accept="image/jpeg"
            @change="loadCoverImage"
          />
        </div>

        <text-input
          v-model="book.title"
          type="text"
          required="true"
          label="Title"
          :value="book.title"
          name="title"
        ></text-input>

        <select-input
          name="author-id"
          v-model="this.book.author_id"
          :items="this.authors"
          required="true"
          label="Author"
        ></select-input>

        <text-input
          v-model="book.publication_year"
          type="number"
          min="0"
          required="true"
          label="Publication Year"
          :value="book.publication_year"
          name="publication_year"
        ></text-input>

        <div class="mb-3">
          <label for="description" class="form-label">Description</label>
          <textarea
            required
            v-model="book.description"
            class="form-control"
            id="description"
            rows="3"
          >
          </textarea>
        </div>

        <div class="mb-3">
          <label for="genres" class="form-label">Genres</label>
          <select
            ref="genres"
            name="genres"
            id="genres"
            required
            size="7"
            v-model="this.book.genre_ids"
            multiple
            class="form-select"
          >
            <option v-for="g in this.genres" :value="g.value" :key="g.value">
              {{ g.text }}
            </option>
          </select>
        </div>

        <hr />

        <div class="float-start">
          <input type="submit" value="Save" class="btn btn-primary me-2" />
          <router-link to="/admin/books" class="btn btn-outline-secondary"
            >Cancel</router-link
          >
        </div>
        <div class="float-end">
          <a
            v-if="this.book.id > 0"
            class="btn btn-danger"
            href="javascript:void(0);"
            @click="confirmDelete(this.book.id)"
            >Delete</a
          >
        </div>
        <div class="clearfix"></div>
      </form-tag>
    </div>
  </div>
</template>

<script>
import Security from "./security.js";
import FormTag from "@/components/forms/FormTag.vue";
import SelectInput from "@/components/forms/SelectInput.vue";
import notie from "notie";
import TextInput from "@/components/forms/TextInput.vue";
import router from "@/router/index.js";

export default {
  name: "BookEdit",
  beforeMount() {
    Security.requireToken();

    // get book for edit if ID > 0
    if (this.$route.params.bookId > 0) {
      fetch(
        process.env.VUE_APP_API_URL +
          "/admin/books/" +
          this.$route.params.bookId,
        Security.requestOptions("")
      )
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            this.$emit("error", data.message);
          } else {
            this.book = data.data;

            let genreArray = [];
            for (let i = 0; i < this.book.genres.length; i++) {
              genreArray.push(this.book.genres[i].id);
            }
            this.book.genre_ids = genreArray;
          }
        });
    }

    // get list of authors for dropdown
    fetch(
      process.env.VUE_APP_API_URL + "/admin/authors/all",
      Security.requestOptions("")
    )
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          this.$emit("error", data.message);
        } else {
          this.authors = data.data;
        }
      });
  },
  data() {
    return {
      book: {
        id: 0,
        title: "",
        author_id: 0,
        publication_year: null,
        description: "",
        cover: "",
        slug: "",
        genres: [],
        genre_ids: [],
      },
      authors: [],
      imgPath: process.env.VUE_APP_IMAGE_URL,
      genres: [
        { value: 1, text: "Science Fiction" },
        { value: 2, text: "Fantasy" },
        { value: 3, text: "Romance" },
        { value: 4, text: "Thriller" },
        { value: 5, text: "Mystery" },
        { value: 6, text: "Horror" },
        { value: 7, text: "Classic" },
      ],
    };
  },
  components: {
    "form-tag": FormTag,
    "text-input": TextInput,
    "select-input": SelectInput,
  },
  methods: {
    submitHandler() {
      const payload = {
        id: this.book.id,
        title: this.book.title,
        author_id: parseInt(this.book.author_id, 10),
        publication_year: parseInt(this.book.publication_year, 10),
        description: this.book.description,
        cover: this.book.cover,
        slug: this.book.slug,
        genre_ids: this.book.genre_ids,
      };

      // console.log(payload);

      fetch(
        `${process.env.VUE_APP_API_URL}/admin/books/save`,
        Security.requestOptions(payload)
      )
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            this.$emit("error", data.message);
          } else {
            this.$emit("success", "Changes saved!");
            router.push("/admin/books");
          }
        })
        .catch((error) => {
          this.$emit("error", error);
        });
    },
    loadCoverImage() {
      // get ref to input
      const file = this.$refs.coverInput.files[0];

      // encode image for JSON using base64
      const reader = new FileReader();

      reader.onloadend = () => {
        const base64string = reader.result
          .replace("data:", "")
          .replace(/^.+,/, "");

        this.book.cover = base64string;
      };

      reader.readAsDataURL(file);
    },
    confirmDelete(id) {
      console.log(id);
      notie.confirm({
        text: "Are you sure?",
        submitText: "Delete",
        submitCallback: () => {
          let payload = { id: id };

          fetch(
            process.env.VUE_APP_API_URL + "/admin/books/delete",
            Security.requestOptions(payload)
          )
            .then((response) => response.json())
            .then((data) => {
              if (data.error) {
                this.$emit("error", data.message);
              } else {
                this.$emit("success", "Changes saved!");
                router.push("/admin/books");
              }
            });
        },
      });
    },
  },
};
</script>

<style scoped>
.book-cover {
  max-width: 10em;
}
</style>
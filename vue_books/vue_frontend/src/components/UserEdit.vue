<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-3">User Edit</h1>
      </div>

      <hr />

      <form-tag
        @userEditEvent="submitHandler"
        name="userform"
        event="userEditEvent"
      >
        <text-input
          v-model="user.first_name"
          type="text"
          required="true"
          label="First Name"
          :value="user.first_name"
          name="first-name"
        ></text-input>

        <text-input
          v-model="user.last_name"
          type="text"
          required="true"
          label="Last Name"
          :value="user.last_name"
          name="last-name"
        ></text-input>

        <text-input
          v-model="user.email"
          type="email"
          required="true"
          label="Email"
          :value="user.email"
          name="email"
        ></text-input>

        <text-input
          v-if="this.user.id === 0"
          v-model="user.password"
          type="password"
          required="true"
          label="Password"
          :value="user.password"
          name="password"
        ></text-input>

        <text-input
          v-else
          v-model="user.password"
          type="password"
          label="Password"
          :value="user.password"
          name="password"
          help="Leave blank to keep existing password"
        ></text-input>

        <div class="form-check">
          <input
            v-model="user.active"
            class="form-check-input"
            type="radio"
            id="user-active"
            :value="1"
          />
          <label class="form-check-label" for="user-active">Active</label>
        </div>

        <div class="form-check">
          <input
            v-model="user.active"
            class="form-check-input"
            type="radio"
            id="user-active-2"
            :value="0"
          />
          <label class="form-check-label" for="user-active-2">Inactive</label>
        </div>

        <hr />

        <div class="float-start">
          <input type="submit" value="Save" class="btn btn-primary me-2" />
          <router-link to="/admin/users" class="btn btn-outline-secondary"
            >Cancel</router-link
          >
        </div>
        <div class="float-end">
          <a
            v-if="
              this.$route.params.userId > 0 &&
              parseInt(String(this.$route.params.userId), 10) !== store.user.id
            "
            class="btn btn-danger"
            href="javascript:void(0);"
            @click="confirmDelete($route.params.userId)"
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
import FormTag from "./forms/FormTag.vue";
import TextInput from "./forms/TextInput.vue";
import notie from "notie";
import { store } from "./store.js";
import router from "./../router/index.js";

export default {
  name: "UserEdit",
  beforeMount() {
    Security.requireToken();

    // user ID > 0 indicates the user already exists
    if (parseInt(String(this.$route.params.userId), 10) > 0) {
      fetch(
        `${process.env.VUE_APP_API_URL}/admin/users/get/${this.$route.params.userId}`,
        Security.requestOptions("")
      )
        .then((response) => response.json())
        .then((response) => {
          if (response.error) {
            this.$emit("error", response.message);
          } else {
            this.user = response;

            // password should be empty for existing users
            this.user.password = "";
          }
        })
        .catch((error) => {
          this.$emit("error", error);
        });
    }
  },
  data() {
    return {
      user: {
        id: 0,
        first_name: "",
        last_name: "",
        email: "",
        password: "",
        active: 0,
      },
      store,
    };
  },
  components: {
    "form-tag": FormTag,
    "text-input": TextInput,
  },
  methods: {
    submitHandler() {
      const payload = {
        id: parseInt(this.$route.params.userId, 10),
        first_name: this.user.first_name,
        last_name: this.user.last_name,
        email: this.user.email,
        password: this.user.password,
        active: this.user.active,
      };

      fetch(
        `${process.env.VUE_APP_API_URL}/admin/users/save`,
        Security.requestOptions(payload)
      )
        .then((response) => response.json())
        .then((response) => {
          if (response.error) {
            this.$emit("error", response.message);
          } else {
            this.$emit("success", "Changes saved!");
            router.push("/admin/users");
          }
        })
        .catch((error) => {
          this.$emit("error", error);
        });
    },
    confirmDelete(id) {
      notie.confirm({
        text: "Are you sure?",
        submitText: "Delete",
        submitCallback: function () {
          const payload = {
            id: id,
          };

          fetch(
            `${process.env.VUE_APP_API_URL}/admin/users/delete`,
            Security.requestOptions(payload)
          )
            .then((response) => response.json())
            .then((response) => {
              if (response.error) {
                this.$emit("error", response.message);
              } else {
                this.$emit("success", "Changes saved!");
                router.push("/admin/users");
              }
            });
        },
      });
    },
  },
};
</script>
<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-5">Login</h1>
        <hr />
        <form-tag>
          <text-input label="Email" type="email" name="email" required="true">
          </text-input>

          <text-input
            label="Password"
            type="password"
            name="password"
            required="true"
          >
          </text-input>

          <hr />
          <input type="submit" class="btn btn-primary" value="Login" />
        </form-tag>
      </div>
    </div>
  </div>
</template>

<script>
import TextInput from "./forms/TextInput.vue";
import FormTag from "./forms/FormTag.vue";

export default {
  name: "UserLogin",
  components: {
    TextInput,
    FormTag,
  },
  mounted() {},
};
</script>
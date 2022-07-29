<template>
  <div class="form-check">
    <input
      class="form-check-input"
      :required="required"
      :name="name"
      :value="value"
      :id="name"
      :checked="checked"
      type="checkbox"
    />
    <label :for="name" class="form-check-label">{{ label }}</label>
  </div>
</template>

<script>
export default {
  name: "CheckboxInput",
  props: ["label", "name", "value", "required", "checked"],
};
</script>
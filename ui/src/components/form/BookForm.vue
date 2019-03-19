<template>
  <div>
    <a-form>
      <a-form-item label="笔记本名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input v-model="book.name" placeholder="请填写笔记本名"></a-input>
      </a-form-item>
      <a-form-item label="颜色" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <swatches v-model="book.color" inline swatch-size="30"/>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import Swatches from 'vue-swatches'
import 'vue-swatches/dist/vue-swatches.min.css'
export default {
  components: {
    Swatches
  },
  computed: {
    book() {
      return this.$store.getters.record;
    }
  },
  watch: {
    book(val) {
      this.$store.dispatch('setRecord', val);
    }
  }
}
</script>

<style scoped>
.vue-swatches__wrapper {
  width: 240px !important;
}
</style>

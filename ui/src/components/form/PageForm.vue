<template>
  <div>
    <a-form>
      <a-form-item label="笔记本" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="page.bookId" placeholder="请选择笔记本"></a-select>
      </a-form-item>
      <a-form-item label="分区" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="page.partId" placeholder="请选择分区"></a-select>
      </a-form-item>
      <a-form-item label="标题" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-textarea v-model="page.title" :autosize="{minRows: 3, maxRows: 6}" placeholder="请填写页面标题"></a-textarea>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
export default {
  computed: {
    page() {
      return this.$store.getters.record;
    }
  },
  watch: {
    page(val) {
      this.$store.dispatch('setRecord', val);
    }
  }
}
</script>

<style scoped>

</style>

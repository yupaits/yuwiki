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
      const page = this.$store.getters.record;
      page.bookId = this.$store.getters.bookId;
      page.partId = this.$store.getters.partId;
      return page;
    },
    partTree() {
      let tree = [];
      if (this.parts && this.parts.length > 0) {
        this.parts.forEach(part => {
          const isLeaf = part.partType === 0;
          tree.push({
            key: part.ID,
            value: part.ID,
            selectable: isLeaf,
            isLeaf: isLeaf,
            scopedSlots: {title: 'part-title'}
          });
        });
      }
      return tree;
    }
  },
  watch: {
    page(val) {
      this.$store.dispatch('setRecord', val);
    }
  },
  data() {
    return {
      books: [],
      parts: []
    }
  },
  created() {
    this.fetchBooks();
    this.fetchBookParts(this.$store.getters.bookId);
  },
  methods: {
    fetchBooks() {
      this.$api.getBooks().then(res => {
        this.books = res.data;
      });
    },
    fetchBookParts(bookId) {
      this.$api.getParts(bookId).then(res => {
        this.parts = res.data;
      });
    },
    handleBookChange(bookId) {
      if (bookId) {
        this.fetchBookParts(bookId);
      } else {
        this.parts = [];
      }
    },
  }
}
</script>

<style scoped>

</style>

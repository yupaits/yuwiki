<template>
  <div>
    <a-form>
      <a-form-item label="笔记本" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="page.bookId" allowClear placeholder="请选择笔记本" @change="handleBookChange">
          <a-select-option v-for="book in books" :key="book.ID" :value="book.ID"><a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="分区" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-tree-select v-model="page.partId" :treeData="partTree" treeNodeLabelProp="label" allowClear placeholder="请选择所属分区"></a-tree-select>
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
      const bookId = this.$store.getters.bookId;
      if (bookId) {
        page.bookId = bookId;
        this.fetchBookParts(bookId);
        const part = this.$store.getters.part;
        if (part) {
          this.$set(page, 'partId', part.partType === 0 ? part.ID : undefined);
        }
      }
      return page;
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
      partTree: []
    }
  },
  created() {
    this.fetchBooks();
  },
  methods: {
    fetchBooks() {
      this.$api.getBooks().then(res => {
        this.books = res.data;
      });
    },
    fetchBookParts(bookId) {
      this.$api.getParts(bookId).then(res => {
        this.partTree = this.toTreeData([], undefined, res.data);
      });
    },
    toTreeData(treeData, prefix, parts) {
      if (parts instanceof Array) {
        parts.forEach(part => {
          const node = {
            title: part.name,
            value: part.ID,
            key: part.ID,
            label: (prefix ? prefix + ' / ' : '') + part.name,
            children: [],
            disabled: part.partType === 1
          }
          this.toTreeData(node.children, node.label, part.SubParts);
          treeData.push(node);
        });
      }
      return treeData;
    },
    handleBookChange(bookId) {
      if (bookId) {
        this.fetchBookParts(bookId);
      } else {
        this.page.partId = undefined;
        this.partTree = [];
      }
    }
  }
}
</script>

<style scoped>

</style>

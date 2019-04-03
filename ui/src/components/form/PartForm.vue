<template>
  <div>
    <a-form>
      <a-form-item label="笔记本" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="part.bookId" allowClear placeholder="请选择笔记本" @change="handleBookChange">
          <a-select-option v-for="book in books" :key="book.ID" :value="book.ID"><a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="分区组" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <a-tree-select v-model="part.parentId" :treeData="partTree" allowClear placeholder="请选择所属分区组"></a-tree-select>
      </a-form-item>
      <a-form-item label="分区名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input v-model="part.name" placeholder="请填写分区名"></a-input>
      </a-form-item>
      <a-form-item label="分区类型" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-radio-group v-model="part.partType" @change="handlePartTypeChange">
          <a-radio :value="0">分区</a-radio>
          <a-radio :value="1">分区组</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="密码保护" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-switch v-model="part.protected"></a-switch>
        <a-input type="password" v-model="part.password" placeholder="请填写保护密码" v-if="part.protected"></a-input>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
export default {
  computed: {
    part() {
      return this.$store.state.record;
    }
  },
  watch: {
    part(val) {
      this.$store.dispatch('setRecord', val);
    }
  },
  data() {
    return {
      books: [],
      partTree: [],
    }
  },
  created() {
    this.fetchBooks();
    const bookId = this.$store.getters.bookId;
    if (bookId) {
      this.part.bookId = bookId;
      this.fetchBookParts(bookId);
    }
  },
  methods: {
    fetchBooks() {
      this.$api.getBooks().then(res => {
        this.books = res.data;
      });
    },
    fetchBookParts(bookId) {
      this.$api.getParts(bookId).then(res => {
        this.partTree = this.toTreeData([], res.data);
      });
    },
    toTreeData(treeData, parts) {
      if (parts instanceof Array) {
        parts.forEach(part => {
          if (part.partType === 1) {
            const node = {
              title: part.name,
              value: part.ID,
              key: part.ID,
              children: []
            };
            this.toTreeData(node.children, part.SubParts);
            treeData.push(node);
          }
        });
      }
      return treeData;
    },
    handleBookChange(bookId) {
      if (bookId) {
        this.fetchBookParts(bookId);
      } else {
        this.part.parentId = undefined;
        this.partTree = [];
      }
    },
    handlePartTypeChange(e) {
      this.$set(this.part, 'partType', e.target.value);
    }
  }
}
</script>

<style scoped>

</style>

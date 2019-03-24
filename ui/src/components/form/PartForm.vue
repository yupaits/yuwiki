<template>
  <div>
    <a-form>
      <a-form-item label="笔记本" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="part.bookId" allowClear placeholder="请选择笔记本" @change="handleBookChange">
          <a-select-option v-for="book in books" :key="book.ID" :value="book.ID"><a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="分区组" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <a-select v-model="part.groupId" :treeData="partTree" allowClear placeholder="请选择所属分区组">
          <template slot="part-title" slot-scope="part">
            <a-icon type="folder"/> {{part.name}}
          </template>
        </a-select>
      </a-form-item>
      <a-form-item label="分区名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input v-model="part.name" placeholder="请填写分区名"></a-input>
      </a-form-item>
      <a-form-item label="分区类型" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-radio-group v-model="partType" @change="handlePartTypeChange">
          <a-radio value="PART">分区</a-radio>
          <a-radio value="GROUP">分区组</a-radio>
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
      return this.$store.getters.record;
    },
    partTree() {
      let tree = [];
      if (this.parts && this.parts.length > 0) {
        this.parts.forEach(part => {
          if (part.partType === 'GROUP') {
            tree.push({
              key: part.ID,
              value: part.ID,
              scopedSlots: {title: 'part-title'}
            });
          }
        })
      }
      return tree;
    }
  },
  watch: {
    part(val) {
      this.partType = val.partType || 'PART';
      this.$store.dispatch('setRecord', val);
    }
  },
  data() {
    return {
      partType: 'PART',
      books: [],
      parts: [],
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
        this.parts = res.data;
      });
    },
    handleBookChange(bookId) {
      this.fetchBookParts(bookId);
    },
    handlePartTypeChange(e) {
      const partType = e.target.value;
      this.partType = partType;
      this.part.partType = partType;
    }
  }
}
</script>

<style scoped>

</style>

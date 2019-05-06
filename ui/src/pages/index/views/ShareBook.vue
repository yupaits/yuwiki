<template>
  <div class="share-page">
    <a-row>
      <a-col :span="8" :offset="8">
        <h3>
          <span><a-icon type="share-alt"/> 共享笔记本</span>
          <span>
            <a-button type="dashed" size="small" class="pull-right" icon="rollback" @click="$router.go(-1)">返回</a-button>
          </span>
        </h3>
        <a-card class="books-card mt-2">
          <a-spin :spinning="loading">
            <div v-if="books.length > 0">
              <div v-for="book in books" :key="book.ID" class="book-item">
                <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
                <span class="pull-right">
                  <a-button size="small" icon="share-alt" @click="shareBook(book)">共享</a-button>
                </span>
              </div>
            </div>
            <div v-else>
              <a-alert message="当前无笔记本，请先创建笔记本再进行共享！"></a-alert>
            </div>
          </a-spin>
        </a-card>
      </a-col>
    </a-row>

    <a-modal :visible.sync="shareVisible" :footer="null" @cancel="shareVisible = false">
      <div slot="title">
        分享笔记本 <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
      </div>
      <div class="user-holder">
        <a-form layout="inline">
          <a-form-item label="选择目标用户">
            <a-input placeholder="请输入用户名">
              <a-icon slot="prefix" type="user"/>
            </a-input>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" icon="check">确认</a-button>
          </a-form-item>
        </a-form>
        <a-divider/>
        <a-form layout="inline">
          <a-form-item label="筛选共享用户">
            <a-input placeholder="请输入关键字">
              <a-icon slot="prefix" type="filter"/>
            </a-input>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script>
export default {
  data() {
    return {
      books: [{}],
      loading: false,
      shareVisible: false,
      book: {},
    }
  },
  created() {
    this.fetchBooks();
  },
  methods: {
    fetchBooks() {
      this.loading = true;
      this.$api.getBooks().then(res => {
        this.books = res.data;
        this.loading = false;
      }).catch(() => {
        this.loading = false;
      });
    },
    shareBook(book) {
      this.book = JSON.parse(JSON.stringify(book));
      this.shareVisible = true;
    },
    handleShareBook() {

    }
  }
}
</script>

<style scoped>
.share-page {
  padding: 16px 24px;
}
.books-card {
  min-height: 500px;
}
.book-item {
  line-height: 32px;
  font-size: 16px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.book-item:hover {
  background: #e6f7ff;
}
.user-holder {
  min-height: 360px;
}
</style>

<template>
  <div class="shared-book-page">
    <a-row :gutter="16">
      <a-col :span="6">
        <h3>
          <span><a-icon type="share-alt"/> 共享给我的笔记本</span>
          <span>
            <a-button type="dashed" size="small" class="pull-right" icon="rollback" @click="$router.go(-1)">返回</a-button>
          </span>
        </h3>
        <div class="mt-3">
          <div>
            <h4 class="text-bold"><a-icon type="book"/> 笔记本</h4>
            <a-select v-model="bookId" class="shared-select" placeholder="请选择笔记本">
              <a-select-option v-for="book in books" :key="book.ID">
                <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
              </a-select-option>
            </a-select>
          </div>
          <div class="mt-1" v-if="bookId">
            <h4 class="text-bold"><a-icon type="folder-open" placeholder="请选择分区"/> 分区</h4>
            <a-select class="shared-select"></a-select>
          </div>
          <div class="mt-1" v-if="partId">
            <h4 class="text-bold"><a-icon type="file-text"/> 页面</h4>
            <a-select class="shared-select" placeholder="请选择页面"></a-select>
          </div>
        </div>
      </a-col>
      <a-col :span="18">
        <mavon-editor :value="this.viewedPage.content" :boxShadow="false" :toolbars="toolbars" :editable="false" defaultOpen="preview" :subfield="false" class="page-preview"></mavon-editor>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import config from '../config'

export default {
  data() {
    return {
      books: [],
      parts: [],
      pages: [],
      bookId: undefined,
      partId: undefined,
      pageId: undefined,
      viewedPage: {},
      toolbars: config.preivew.toolbars
    }
  },
  created() {
    this.fetchSharedBooks();
  },
  methods: {
    fetchSharedBooks() {
      this.$api.getSharedBooks().then(res => {
        this.books = res.data;
      });
    },
    fetchParts(bookId) {
      this.$api.getParts(bookId).then(res => {
        this.parts = res.data;
      });
    },
    fetchPages(partId) {
      this.$api.getPages(partId).then(res => {
        this.pages = res.data;
      });
    },
  }
}
</script>

<style scoped>
.shared-book-page {
  padding: 16px 24px;
}
.shared-select {
  width: 80%;
}
.page-preview {
  height: calc(100vh - 32px);
  z-index: 0;
}
</style>

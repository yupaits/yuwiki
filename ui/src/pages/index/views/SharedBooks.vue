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
            <a-select v-model="bookId" allowClear class="shared-select" placeholder="请选择笔记本" @change="handleBookChange">
              <a-select-option v-for="book in books" :key="book.ID">
                <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
              </a-select-option>
            </a-select>
          </div>
          <div class="mt-1" v-if="bookId">
            <h4 class="text-bold"><a-icon type="folder-open" placeholder="请选择分区"/> 分区</h4>
            <a-tree-select v-model="partId" :treeData="partTree" treeNodeLabelProp="label" allowClear class="shared-select" placeholder="请选择所属分区" @change="handlePartChange"></a-tree-select>
          </div>
          <div class="mt-1" v-if="partId">
            <h4 class="text-bold"><a-icon type="file-text"/> 页面</h4>
            <a-select v-model="pageId" allowClear class="shared-select" placeholder="请选择页面" @change="handlePageChange">
              <a-select-option v-for="page in pages" :key="page.ID">
                <a-icon type="file-text"/> {{page.title}}
              </a-select-option>
            </a-select>
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
      partTree: [],
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
      this.$api.getSharedParts(bookId).then(res => {
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
    fetchPages(partId) {
      this.$api.getSharedPages(this.bookId, partId).then(res => {
        this.pages = res.data;
      });
    },
    viewPage(pageId) {
      this.$api.getSharedPage(this.bookId, this.partId, pageId).then(res => {
        this.viewedPage = res.data;
      });
    },
    handleBookChange(bookId) {
      if (bookId) {
        this.fetchParts(bookId);
      } else {
        this.partId = undefined;
        this.partTree = [];
        this.viewedPage = {};
      }
    },
    handlePartChange(partId) {
      if (partId) {
        this.fetchPages(partId);
      } else {
        this.pageId = undefined;
        this.viewedPage = {};
      }
    },
    handlePageChange(pageId) {
      if (pageId) {
        this.viewPage(pageId);
      } else {
        this.viewedPage = {};
      }
    }
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

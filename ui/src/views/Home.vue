<template>
  <div>
    <a-layout id="components-layout-demo-top" class="layout">
      <a-layout-header>
        <span class="logo">
          <img src="favicon.ico" alt="YuWIki">
          <span class="ml-1">知识库</span>
        </span>
        <span class="pull-right">
          <a-button :icon="menuVisible ? 'menu-fold' : 'menu-unfold'" @click="menuVisible = !menuVisible"></a-button>
          <a-input-search class="search-input ml-1"></a-input-search>
          <a-dropdown>
            <a-menu slot="overlay" @click="handleCreate">
              <a-menu-item key="book"><a-icon type="book"/>笔记本</a-menu-item>
              <a-menu-item key="part"><a-icon type="folder-open"/>分区</a-menu-item>
              <a-menu-item key="page"><a-icon type="file-text"/>页面</a-menu-item>
            </a-menu>
            <a-button icon="plus" class="ml-1">
              新建 <a-icon type="down"/>
            </a-button>
          </a-dropdown>
          <a-dropdown placement="bottomRight">
            <a-menu slot="overlay" @click="handleUserOpt">
              <a-menu-item key="share-book"><a-icon type="share-alt"/>共享笔记本</a-menu-item>
              <a-menu-item key="modify-passwd"><a-icon type="key"/>修改密码</a-menu-item>
              <a-menu-item key="logout"><a-icon type="logout"/>注销登录</a-menu-item>
            </a-menu>
            <span class="user-holder ml-1">
              <a-avatar size="small" shape="square" icon="user"></a-avatar>
              <span class="ml-1">yupaits</span> <a-icon type="down"/>
            </span>
          </a-dropdown>
        </span>
      </a-layout-header>
      <a-layout-content class="layout-content">
        <a-row :gutter="16">
          <div v-if="menuVisible">
            <a-col :span="4">
              <div class="holder">
                <h3 class="text-title text-bold holder-header">
                  <a-icon type="book"/> 笔记本
                  <a-button size="small" icon="sync" @click="fetchBooks"></a-button>
                  <span class="pull-right" v-if="$store.getters.bookId">
                    <a-button size="small" icon="edit" class="mr-1" @click="editBook"></a-button>
                    <a-popconfirm title="确定删除此笔记本吗？" placement="right" @confirm="handleDeleteBook">
                      <a-button size="small" icon="delete"></a-button>
                    </a-popconfirm>
                  </span>
                </h3>
                <a-spin :spinning="loading.books" class="list">
                  <div v-for="book in books" :key="book.ID" class="book-item" :class="{'active': $store.getters.bookId === book.ID}" @click="selectBook(book.ID)">
                    <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
                  </div>
                </a-spin>
              </div>
            </a-col>
            <a-col :span="4">
              <div class="holder">
                <h3 class="text-title text-bold holder-header">
                  <a-icon type="folder-open"/> 分区
                  <span v-if="$store.getters.bookId">
                    <a-button size="small" icon="sync" class="ml-1" @click="fetchParts($store.getters.bookId)"></a-button>
                    <span class="pull-right" v-if="$store.getters.partId">
                      <a-button size="small" icon="edit" class="mr-1" @click="editPart"></a-button>
                      <a-popconfirm title="确定删除此分区吗？" placement="right" @confirm="handleDeletePart">
                        <a-button size="small" icon="delete"></a-button>
                      </a-popconfirm>
                    </span>
                  </span>
                </h3>
                <a-spin :spinning="loading.parts" class="list">
                  <part-tree :parts="parts"/>
                </a-spin>
              </div>
            </a-col>
          </div>
          <div v-else>
            <a-col :span="1">
              <div class="holder fold-holder" @click="menuVisible = true">
                <h3 class="text-title text-bold holder-header">
                  <a-icon type="book"/> 笔记本
                  <a-divider></a-divider>
                  <a-icon type="folder-open"/> 分区
                </h3>
              </div>
            </a-col>
          </div>
          <a-col :span="6">
            <div class="holder">
              <h3 class="text-title text-bold holder-header">
                <a-icon type="file-text"/> 页面
                <span v-if="$store.getters.partId">
                  <a-button size="small" icon="sync" class="ml-1" @click="fetchPages($store.getters.partId)"></a-button>
                  <span class="pull-right" v-if="$store.getters.pageId">
                    <a-button size="small" icon="form" class="mr-3" @click="toEditor"></a-button>
                    <a-button size="small" icon="edit" class="mr-1" @click="editPage"></a-button>
                    <a-popconfirm title="确定删除此页面吗？" placement="right" @confirm="handleDeletePage">
                      <a-button size="small" icon="delete"></a-button>
                    </a-popconfirm>
                  </span>
                </span>
              </h3>
              <a-spin :spinning="loading.pages" class="list">
                <div v-for="page in pages" :key="page.ID" class="page-item" :class="{'active': $store.getters.pageId === page.ID}" @click="selectPage(page.ID)">
                  <a-icon type="file-text"/> {{page.title}}
                </div>
              </a-spin>
            </div>
          </a-col>
          <a-col :span="menuVisible ? 10 : 17">
            <div class="holder preview-holder">
              <mavon-editor :value="this.viewedPage.content" :toolbarsFlag="false" :editable="false" defaultOpen="preview" :subfield="false" class="page-preview"></mavon-editor>
            </div>
          </a-col>
        </a-row>
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        <b>YuWiki</b> ©2019 <b><a href="https://github.com/YupaiTS" target="_blank">YupaiTS</a></b> 版权所有
      </a-layout-footer>
    </a-layout>

    <a-modal :visible.sync="modal.visible" @cancel="closeModal" @ok="modal.ok">
      <template slot="title">
        <a-icon :type="modalType[modal.type] ? modalType[modal.type].icon : 'question'"/> {{modalTitle}}
      </template>
      <component :is="form"></component>
    </a-modal>
  </div>
</template>

<script>
  import BookForm from '../components/form/BookForm'
  import PartForm from '../components/form/PartForm'
  import PageForm from '../components/form/PageForm'
  import PartTree from "../components/PartTree"
import { unlink } from 'fs';
  export default {
    components: {
      PartTree
    },
    data() {
      return {
        books: [],
        parts: [],
        pages: [],
        viewedPage: {},
        loading: {
          books: false,
          parts: false,
          pages: false,
          pageView: false
        },
        menuVisible: true,
        modalVisible: false,
        modal: {
          type: undefined,
          key: undefined,
          visible: false,
          ok: () => {
          }
        },
        modalType: {
          create: {label: '新建', icon: 'plus'},
          modify: {label: '修改', icon: 'edit'}
        },
        title: {
          book: '笔记本',
          part: '分区',
          page: '页面',
        },
        createHandlers: {
          book: this.handleAddBook,
          part: this.handleAddPart,
          page: this.handleAddPage,
        },
        forms: {
          book: BookForm,
          part: PartForm,
          page: PageForm,
        },
        default: {
          book: {},
          part: {
            partType: 0,
            protected: false
          },
          page: {}
        }
      }
    },
    computed: {
      modalTitle() {
        const type = this.modalType[this.modal.type];
        return (type ? type.label : '') + this.title[this.modal.key];
      },
      selectedBook() {
        return this.books.filter(book => book.ID === this.$store.getters.bookId)[0] || {};
      },
      form() {
        return this.forms[this.modal.key];
      }
    },
    created() {
      this.fetchBooks();
      this.$eventBus.$on('selectPart', this.selectPart);
      if (this.$store.getters.pageId) {
        this.menuVisible = false;
        this.fetchParts(this.$store.getters.bookId);
        this.fetchPages(this.$store.getters.partId);
        this.viewPage();
      }
    },
    methods: {
      fetchBooks() {
        this.loading.books = true;
        this.$api.getBooks().then(res => {
          this.books = res.data;
          this.loading.books = false;
        }).catch(() => {
          this.loading.books = false;
        });
      },
      fetchParts(bookId) {
        this.loading.parts = true;
        this.$api.getParts(bookId).then(res => {
          this.parts = res.data;
          this.loading.parts = false;
        }).catch(() => {
          this.loading.parts = false;
        });
      },
      fetchPages(partId) {
        this.loading.pages = true;
        this.$api.getPages(partId).then(res => {
          this.pages = res.data;
          this.loading.pages = false;
        }).catch(() => {
          this.loading.pages = false;
        });
      },
      viewPage() {
        this.loading.pageView = true;
        this.$api.viewPage(this.$store.getters.pageId).then(res => {
          this.viewedPage = res.data;
          this.loading.pageView = false;
        }).catch(() => {
          this.loading.pageView = false;
        });
      },
      handleUserOpt({key}) {
        alert(key);
      },
      showModal(type, key) {
        this.modal = {
          type: type,
          key: key,
          visible: true,
          ok: this.createHandlers[key]
        };
      },
      handleCreate({key}) {
        this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.default[key])));
        this.showModal('create', key);
      },
      closeModal() {
        this.modal.visible = false;
      },
      selectBook(bookId) {
        this.$store.dispatch('setBookId', bookId);
        this.$store.dispatch('setPartId', undefined);
        this.$store.dispatch('setPageId', undefined);
        this.pages = [];
        this.viewedPage = {};
        if (bookId) {
          this.fetchParts(bookId);
        }
      },
      selectPart(partId) {
        this.$store.dispatch('setPageId', undefined);
        this.viewedPage = {};
        if (partId) {
          this.fetchPages(partId);
        }
      },
      selectPage(pageId) {
        this.$store.dispatch('setPageId', pageId);
        if (pageId) {
          this.viewPage();
        }
      },
      editBook() {
        this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.selectedBook)));
        this.modal = {
          type: 'modify',
          key: 'book',
          visible: true,
          ok: this.handleEditBook
        };
      },
      editPart() {
        this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.$store.getters.part)));
        this.modal = {
          type: 'modify',
          key: 'part',
          visible: true,
          ok: this.handleEditPart
        };
      },
      editPage() {
        this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.viewedPage)));
        this.modal = {
          type: 'modify',
          key: 'page',
          visible: true,
          ok: this.handleEditPage
        };
      },
      toEditor() {
        this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.viewedPage)));
        this.$router.push('/page/edit');
      },
      handleAddBook() {
        this.$api.addBook(this.$store.getters.record).then(() => {
          this.$message.success(this.$messages.result.createSuccess);
          this.fetchBooks();
          this.closeModal();
        });
      },
      handleAddPart() {
        this.$api.addPart(this.$store.getters.record).then(() => {
          this.$message.success(this.$messages.result.createSuccess);
          if (this.$store.getters.bookId) {
            this.fetchParts(this.$store.getters.bookId);
          }
          this.closeModal();
        });
      },
      handleAddPage() {
        this.$api.addPage(this.$store.getters.record).then(() => {
          this.$message.success(this.$messages.result.createSuccess);
          if (this.$store.getters.partId) {
            this.fetchPages(this.$store.getters.partId);
          }
          this.closeModal();
        });
      },
      handleEditBook() {
        this.$api.editBook(this.$store.getters.record).then(() => {
          this.$message.success(this.$messages.result.updateSuccess);
          this.fetchBooks();
          this.closeModal();
        });
      },
      handleEditPart() {
        const part = this.$store.getters.record;
        this.$api.editPart(part).then(() => {
          this.$message.success(this.$messages.result.updateSuccess);
          this.fetchParts(this.$store.getters.bookId);
          this.$store.dispatch('setPart', part);
          this.closeModal();
        });
      },
      handleEditPage() {
        this.$api.updatePage(this.$store.getters.record).then(() => {
          this.$message.success(this.$messages.result.updateSuccess);
          this.fetchPages(this.$store.getters.partId);
          this.closeModal();
        });
      },
      handleDeleteBook() {
        this.$api.deleteBook(this.$store.getters.bookId).then(() => {
          this.$message.success(this.$messages.result.deleteSuccess);
          this.fetchBooks();
          this.selectBook(undefined);
        });
      },
      handleDeletePart() {
        this.$api.deletePart(this.$store.getters.partId).then(() => {
          this.$message.success(this.$messages.result.deleteSuccess);
          this.fetchParts(this.$store.getters.bookId);
          this.selectPart(undefined);
        });
      },
      handleDeletePage() {
        this.$api.deletePage(this.$store.getters.pageId).then(() => {
          this.$message.success(this.$messages.result.deleteSuccess);
          this.fetchPages(this.$store.getters.partId);
          this.selectPage(undefined);
        });
      }
    }
  }
</script>

<style scoped>
::-webkit-scrollbar
{
  width: 6px;
  background-color: #f5f5f5;
}
::-webkit-scrollbar-track
{
  border-radius: 10px;
  background-color: #f5f5f5;
}
::-webkit-scrollbar-thumb
{
  border-radius: 10px;
  background-color: #bfbfbf;
}
.ant-layout-header {
  background: #f0f2f5;
}
.logo {
  font-size: 24px;
  line-height: 64px;
}
.search-input {
  width: 320px;
}
.user-holder {
  padding: 8px 16px;
  border: 1px solid #f5f5f5;
  border-radius: 4px;
  background: #fff;
}
.layout-content {
  height: calc(100vh - 133px);
  padding: 0 50px;
}
.holder {
  height: calc(100vh - 133px);
  padding: 8px 16px;
  border: 1px solid #f5f5f5;
  border-radius: 4px;
  background: #fff;
}
.fold-holder {
  cursor: pointer;
}
.holder-header {
  line-height: 28px;
}
.list {
  height: calc(100vh - 210px);
  overflow-x: hidden;
  overflow-y: auto;
}
.book-item,.page-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.book-item:hover,.page-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.book-item.active,.page-item.active {
  background: #91d5ff;
  font-weight: bold;
  color: #262626;
}
.preview-holder {
  padding: 2px;
}
.page-preview {
  z-index: 0;
  height: 100%;
}
</style>

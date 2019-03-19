<template>
  <div>
    <a-layout id="components-layout-demo-top" class="layout">
      <a-layout-header>
        <span class="logo">
          <img src="favicon.ico" alt="YuWIki">
          <span class="ml-1">知识库</span>
        </span>
        <span class="pull-right">
          <a-input-search class="search-input"></a-input-search>
          <a-dropdown>
            <a-menu slot="overlay" @click="handleCreate">
              <a-menu-item key="book"><a-icon type="book"/>笔记本</a-menu-item>
              <a-menu-item key="part"><a-icon type="folder"/>分区</a-menu-item>
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
              <a-avatar size="small" shape="square" icon="user" src="https://avatars1.githubusercontent.com/u/12194490?s=40&v=4"></a-avatar>
              <span class="ml-1">yupaits</span> <a-icon type="down"/>
            </span>
          </a-dropdown>
        </span>
      </a-layout-header>
      <a-layout-content class="layout-content">
        <a-row :gutter="16">
          <a-col :span="4">
            <div class="holder">
              <h3 class="text-title text-bold holder-header">
                <a-icon type="book"/> 笔记本
                <a-button size="small" icon="sync" class="ml-1" @click="fetchBooks"></a-button>
                <span class="pull-right" v-if="bookId">
                  <a-button size="small" icon="edit" class="mr-1" @click="editBook"></a-button>
                  <a-popconfirm title="确定删除此笔记本吗？" placement="right" @confirm="handleDeleteBook">
                    <a-button size="small" icon="delete"></a-button>
                  </a-popconfirm>
                </span>
              </h3>
              <a-spin :spinning="loading.books" class="list">
                <div v-for="book in books" :key="book.id" class="book-item" :class="{'active': bookId === book.id}" @click="selectBook(book.id)">
                  <span :style="{color: book.color}"><a-icon type="book"/></span> {{book.name}}
                </div>
              </a-spin>
            </div>
          </a-col>
          <a-col :span="4">
            <div class="holder">
              <h3 class="text-title text-bold holder-header">
                <a-icon type="folder"/> 分区
                <span v-if="partId">
                  <a-button size="small" icon="sync" class="ml-1" @click="fetchParts"></a-button>
                  <span class="pull-right">
                    <a-button size="small" icon="edit" class="mr-1" @click="editPart"></a-button>
                    <a-popconfirm title="确定删除此分区吗？" placement="right" @confirm="handleDeletePart">
                      <a-button size="small" icon="delete"></a-button>
                    </a-popconfirm>
                  </span>
                </span>
              </h3>
            </div>
          </a-col>
          <a-col :span="6">
            <div class="holder">
              <h3 class="text-title text-bold holder-header">
                <a-icon type="file-text"/> 页面
                <span v-if="pageId">
                  <a-button size="small" icon="sync" class="ml-1" @click="fetchParts"></a-button>
                  <span class="pull-right">
                    <a-button size="small" icon="edit" class="mr-1" @click="editPart"></a-button>
                    <a-popconfirm title="确定删除此分区吗？" placement="right" @confirm="handleDeletePart">
                      <a-button size="small" icon="delete"></a-button>
                    </a-popconfirm>
                  </span>
                </span>
              </h3>
            </div>
          </a-col>
          <a-col :span="10">
            <div class="holder">
            </div>
          </a-col>
        </a-row>
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        <b>YuWiki</b> ©2019 <b><a href="https://github.com/YupaiTS" target="_blank">YupaiTS</a></b> 版权所有
      </a-layout-footer>
    </a-layout>

    <a-modal :visible.sync="modal.visible" @cancel="handleClose" @ok="modal.ok">
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
import ModifyPasswordForm from '../components/form/ModifyPasswordForm'
import PageEditor from '../components/PageEditor'
export default {
  components: {
    PageEditor
  },
  data() {
    return {
      books: [{id: 1, name: '读后感'}],
      parts: [],
      pages: [],
      bookId: undefined,
      partId: undefined,
      pageId: undefined,
      viewedPage: {},
      loading: {
        books: false,
        parts: false,
        pages: false,
        pageView: false
      },
      modalVisible: false,
      modal: {
        type: undefined,
        key: undefined,
        visible: false,
        ok: () => {}
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
      }
    }
  },
  computed: {
    modalTitle() {
      const type = this.modalType[this.modal.type];
      return (type ? type.label : '') + this.title[this.modal.key];
    },
    selectedBook() {
      return this.books.filter(book => book.id === this.bookId)[0] || {};
    },
    selectedPart() {
      return this.parts.filter(part => part.id === this.partId)[0] || {};
    },
    form() {
      return this.forms[this.modal.key];
    }
  },
  created() {
    this.fetchBooks();
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
    fetchParts() {
      this.loading.parts = true;
      this.$api.getParts(this.bookId).then(res => {
        this.parts = res.data;
        this.loading.parts = false;
      }).catch(() => {
        this.loading.parts = false;
      });
    },
    fetchPages() {
      this.loading.pages = true;
      this.$api.getPages(this.partId).then(res => {
        this.pages = res.data;
        this.loading.pages = false;
      }).catch(() => {
        this.loading.pages = false;
      });
    },
    viewPage() {
      this.loading.pageView = true;
      this.$api.viewPage(this.pageId).then(res => {
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
      this.$store.dispatch('setRecord', {});
      this.showModal('create', key);
    },
    handleClose() {
      this.modal.visible = false;
    },
    selectBook(bookId) {
      this.bookId = bookId;
      this.partId = undefined;
      this.pageId = undefined;
      this.viewedPage = {};
      this.fetchParts();
    },
    selectPart(partId) {
      this.partId = partId;
      this.pageId = undefined;
      this.viewedPage = {};
      this.fetchPages();
    },
    selectPage(pageId) {
      this.pageId = pageId;
      this.viewPage();
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
      this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.selectedPart)));
      this.modal = {
        type: 'modify',
        key: 'part',
        visible: true,
        ok: this.handleEditPart
      };
    },
    editPage() {
      this.$store.dispatch('setRecord', JSON.parse(JSON.stringify(this.viewedPage)));
    },
    handleAddBook() {
      this.$api.addBook(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.createSuccess);
        this.fetchBooks();
      });
    },
    handleAddPart() {
      this.$api.addPart(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.createSuccess);
        this.fetchParts();
      });
    },
    handleAddPage() {
      this.$api.addPage(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.createSuccess);
        this.fetchPages();
      });
    },
    handleEditBook() {
      this.$api.editBook(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.updateSuccess);
        this.fetchBooks();
      });
    },
    handleEditPart() {
      this.$api.editPart(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.updateSuccess);
        this.fetchParts();
      });
    },
    handleEditPage() {
      this.$api.editPage(this.$store.getters.record).then(res => {
        this.$message.success(this.$messages.updateSuccess);
        this.fetchPages();
      });
    },
    handleDeleteBook() {
      this.$api.deleteBook(this.bookId).then(() => {
        this.$message.success(this.$messages.deleteSuccess);
        this.fetchBooks();
        this.selectBook(undefined);
      });
    },
    handleDeletePart() {
      this.$api.deletePart(this.partId).then(() => {
        this.$message.success(this.$messages.deleteSuccess);
        this.fetchParts();
        this.selectPart(undefined);
      });
    },
    handleDeletePage() {
      this.$api.deletePage(this.pageId).then(() => {
        this.$message.success(this.$messages.deleteSuccess);
        this.fetchPages();
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
.holder-header {
  line-height: 28px;
}
.list {
  height: calc(100vh - 210px);
  overflow-x: hidden;
  overflow-y: auto;
}
.book-item,.part-item,.page-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
}
.book-item:hover,.part-item:hover,.page-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.book-item.active,.part-item.active,.page-item.active {
  background: #91d5ff;
  font-weight: bold;
  color: #262626;
}
</style>

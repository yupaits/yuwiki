<template>
  <div class="template-page">
    <a-row>
      <a-col :span="8" :offset="8">
        <h3>
          <span><a-icon type="file-markdown"/> 页面模板</span>
          <a-button size="small" icon="plus" title="添加模板" class="ml-2" @click="addTemplate"></a-button>
          <span>
            <a-button type="dashed" size="small" class="pull-right" icon="rollback" @click="$router.go(-1)">返回</a-button>
          </span>
        </h3>
        <a-card class="template-card mt-2">
          <a-spin :spinning="loading">
            <div v-if="templates.length > 0">
              <div v-for="template in templates" :key="template.ID" class="template-item">
                <a-icon type="file-markdown"/> {{template.name}}
                <span class="pull-right">
                  <a-button size="small" icon="eye" title="查看" @click="showTemplate(template.ID)"></a-button>
                  <a-button size="small" icon="edit" title="修改" class="ml-1" @click="modifyTemplate(template.ID)"></a-button>
                  <a-popconfirm title="确定删除此模板吗？" placement="right" @confirm="handleDeleteTemplate(template.ID)">
                    <a-button size="small" icon="delete" title="删除" class="ml-1"></a-button>
                  </a-popconfirm>
                </span>
              </div>
            </div>
            <div v-else>
              <a-alert message="当前无页面模板，可以点击上方的添加按钮创建一个模板！"></a-alert>
            </div>
          </a-spin>
        </a-card>
      </a-col>
    </a-row>

    <a-modal :title="modal.title" :visible.sync="modal.visible" :maskClosable="false" @ok="modal.ok" @cancel="modal.visible = false" :width="1200">
      <a-form>
        <a-form-item label="模板标题">
          <a-input v-model="template.name" placeholder="请填写模板标题"></a-input>
        </a-form-item>
        <a-form-item label="模板格式">
          <mavon-editor :value="template.content" :boxShadow="false" :toolbars="editToolbars" ref="editor" class="editor"></mavon-editor>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal :title="template.name" :visible.sync="showVisible" :maskClosable="false" :footer="null" @cancel="showVisible = false" :width="1200">
      <mavon-editor :value="template.content" :boxShadow="false" :toolbars="previewToolbars" :editable="false" defaultOpen="preview" :subfield="false" class="template-preview"></mavon-editor>
    </a-modal>
  </div>
</template>

<script>
import config from '../config'
export default {
  data() {
    return {
      templates: [],
      loading: false,
      modal: {
        visible: false,
        title: '',
        ok: () => {}
      },
      showVisible: false,
      template: {},
      editToolbars: config.editor.toolbars,
      previewToolbars: config.preview.toolbars,
    }
  },
  created() {
    this.fetchTemplates();
  },
  methods: {
    fetchTemplates() {
      this.loading = true;
      this.$api.getTemplates().then(res => {
        this.templates = res.data || [];
        this.loading = false;
      }).catch(() => {
        this.loading = false;
      });
    },
    addTemplate() {
      this.template = {};
      this.modal = {
        visible: true,
        title: '添加模板',
        ok: this.handleAddTemplate
      };
    },
    showTemplate(templateId) {
      this.$api.getTemplate(templateId).then(res => {
        this.template = res.data;
        this.showVisible = true;
      })
    },
    modifyTemplate(templateId) {
      this.$api.getTemplate(templateId).then(res => {
        this.template = res.data;
        this.modal = {
          visible: true,
          title: '修改模板',
          ok: this.handleModifyTemplate
        };
      });
    },
    handleDeleteTemplate(templateId) {
      this.$api.deleteTemplate(templateId).then(() => {
        this.$message.success(this.$messages.result.deleteSuccess);
        this.fetchTemplates();
      });
    },
    handleAddTemplate() {
      this.$api.addTemplate(this.template).then(() => {
        this.$message.success(this.$messages.result.createSuccess);
        this.modal.visible = false;
        this.fetchTemplates();
      });
    },
    handleModifyTemplate() {
      this.template.content = this.$refs.editor.d_value;
      this.$api.modifyTemplate(this.template).then(() => {
        this.$message.success(this.$messages.result.updateSuccess);
        this.modal.visible = false;
        this.fetchTemplates();
      });
    }
  }
}
</script>

<style scoped>
.template-page {
  padding: 16px 24px;
}
.template-card {
  min-height: 500px;
}
.template-item {
  line-height: 32px;
  font-size: 16px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.template-item:hover {
  background: #e6f7ff;
}
.editor {
  z-index: 0;
  height: 400px;
}
.template-preview {
  z-index: 0;
  min-height: 500px;
}
</style>

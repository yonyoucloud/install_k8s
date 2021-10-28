<template>
  <editor v-model="editorText" :lang="lang" :width="width" :height="height" :options="editorOptions" @init="init"></editor>
</template>

<script>
import { dump as toYaml, load as fromYaml } from 'js-yaml'

export default {
  name: 'AceEditor',
  components: { 
    editor: require('vue2-ace-editor')
  },
  props: {
    text: {
      type: Object,
      default: () => {},
      required: true
    },
    lang: {
      type: String,
      default: 'yaml',
      required: true
    },
    width: {
      type: String
    },
    height: {
      type: String,
      default: '300px'
    },
    options: {
      type: Function,
      default: () => {}
    }
  },
  watch: {
    text: {
      deep: true,
      handler(val) {
        switch (this.lang) {
          case 'yaml':
            this.editorText = toYaml(val)
            break
          case 'json':
            this.editorText = this.toRawJSON(val)
            break
        }
      }
    },
    editorText: {
      deep: true,
      handler(val) {
        let data = {}
        switch (this.lang) {
          case 'yaml':
            data = fromYaml(val)
            break
          case 'json':
            data = JSON.parse(val)
            break
        }
        this.$emit('handleTextChange', this.lang, data)
      }
    }
  },
  data() {
    return {
      editorText: '',
      editorOptions: {
        showPrintMargin: false,
        highlightActiveLine: true,
        tabSize: 2,
        wrap: true,
        fontSize: 14,
        fontFamily: `'Roboto Mono Regular', monospace`,
      }
    }
  },
  methods: {
    init() {
      require('brace')
      require('brace/mode/json')                
      require('brace/mode/yaml')    //language
      require('brace/theme/idle_fingers')
      require('brace/theme/textmate')
      require('brace/worker/json') //snippet
      require('brace/theme/chrome')
    },
    toRawJSON(object) {
      return JSON.stringify(object, null, '\t');
    },
  }
}
</script>
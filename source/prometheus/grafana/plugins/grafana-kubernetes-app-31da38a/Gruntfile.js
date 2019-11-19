module.exports = function(grunt) {
  require('load-grunt-tasks')(grunt);

  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-typescript');
  grunt.loadNpmTasks('grunt-contrib-watch');


  grunt.initConfig({
    clean: ['dist'],

    copy: {
      dist_js: {
        expand: true,
        cwd: 'src',
        src: ['**/*.ts', '**/*.d.ts'],
        dest: 'dist'
      },
      dist_html: {
        expand: true,
        cwd: 'src',
        src: ['**/*.html', '**/*.json'],
        dest: 'dist'
      },
      dist_css: {
        expand: true,
        flatten: true,
        cwd: 'src/css',
        src: ['*.css'],
        dest: 'dist/css/'
      },
      dist_img: {
        expand: true,
        flatten: true,
        cwd: 'src/img',
        src: ['*.*'],
        dest: 'dist/img/'
      },
      dist_statics: {
        expand: true,
        flatten: true,
        src: ['src/plugin.json', 'LICENSE', 'README.md', 'src/query_help.md'],
        dest: 'dist/'
      }
    },

    typescript: {
      build: {
        src: ['dist/**/*.ts', '!**/*.d.ts'],
        dest: 'dist',
        options: {
          module: 'system',
          target: 'es5',
          rootDir: 'dist/',
          declaration: true,
          emitDecoratorMetadata: true,
          experimentalDecorators: true,
          sourceMap: true,
          noImplicitAny: false,
        }
      }
    },

    sass: {
      options: {
        sourceMap: true
      },
      dist: {
        files: {
          "dist/css/kubernetes.dark.css": "src/sass/kubernetes.dark.scss",
          "dist/css/kubernetes.light.css": "src/sass/kubernetes.light.scss",
        }
      }
    },

    watch: {
      files: ['src/**/*.ts', 'src/**/*.html', 'src/**/*.css', 'src/img/*.*', 'src/plugin.json', 'README.md', 'src/query_help.md'],
      tasks: ['default'],
      options: {
        debounceDelay: 250,
      },
    }
  });
  
  grunt.registerTask('default', [
    'clean',
    'copy:dist_js',
    'sass',
    'typescript:build',
    'copy:dist_html',
    'copy:dist_css',
    'copy:dist_img',
    'copy:dist_statics'
  ]);
};
/* jshint node: true */

module.exports = function(grunt) {
  "use strict";

  // Project configuration.
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),

    clean: {
      interim: ['static/js/interim'],
      dist: ['static/js/build']
    },

    react: {
      jsxhint: {
        files: [{
          expand: true,
          cwd: 'static/js/src/components',
          src: ['*.js'],
          dest: 'static/js/interim/jsxhint/',
        }]
      }
    },

    browserify: {
      build: {
        options: {
          transform: [['reactify', {harmony: true}], 'envify'],
          bundleOptions: {
            debug: true,
            detectGlobals: false
          },
        },
        src: ['static/js/src/**/*.js', 'static/js/src/*.js'],
        dest: 'static/js/build/tchaik.js'
      },
    },

    jshint: {
      gruntfile: {
        src: ['Gruntfile.js']
      },
      packagejson: {
        src: ['package.json']
      },
      jsx: {
        files: {
          src: ['static/js/interim/jsxhint/*.js']
        },
        options: {
          esnext: true,
          browser: true,
          devel: true,
          jquery: true,
          curly: true,
          undef: true,
          unused: true,
          node: true,
          newcap: false
        }
      },
      js: {
        files: {
          src: ['static/js/src/*.js', 'static/js/src/actions/*.js', 'static/js/src/stores/*.js', 'static/js/src/dispatcher/*.js', 'static/js/src/constants/*.js', 'static/js/src/api/*.js']
        },
        options: {
          esnext: true,
          browser: true,
          devel: true,
          curly: true,
          undef: true,
          unused: true,
          node: true,
          newcap: false
        }
      },
      tchaik: {
        files: {
          src: ['static/js/build/tchaik.js']
        },
        options: {
          browser: true,
          devel: true,
          jquery: true,
          undef: true,
          unused: true,
          node: true,
          newcap: false
        }
      }
    },

    compass: {
      build: {
        options: {
          httpPath: '/',
          cssPath: 'static/css',
          sassPath: 'static/sass',
          imagesPath: 'static/img',
          javascriptsPath: 'static/js',

          noLineComments: true
        }
      }
    },

    uglify: {
      dist: {
        files: {
          'static/js/build/tchaik.min.js': ['static/js/build/tchaik.js']
        }
      }
    },

    watch: {
      options: {
        spawn: true
      },
      gruntfile: {
        files: 'Gruntfile.js',
        tasks: ['jshint:gruntfile'],
      },
      jsx: {
        files: ['static/js/src/components/*.js'],
        tasks: ['react', 'jshint:jsx', 'clean:interim', 'browserify'],
      },
      js: {
        files: ['static/js/src/*.js', 'static/js/src/actions/*.js', 'static/js/src/stores/*.js', 'static/js/src/dispatcher/*.js', 'static/js/src/constants/*.js', 'static/js/src/api/*.js'],
        tasks: ['jshint:js', 'browserify'],
      },
      sass: {
        files: ['static/sass/*.scss'],
        tasks: ['compass']
      },
      packagejson: {
        files: ['package.json'],
        tasks: ['jshint:packagejson']
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-watch');

  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-react');
  grunt.loadNpmTasks('grunt-browserify');
  grunt.loadNpmTasks('grunt-contrib-jshint');

  grunt.loadNpmTasks('grunt-contrib-uglify');

  grunt.loadNpmTasks('grunt-contrib-compass');

  grunt.registerTask('default', ['browserify']);
};

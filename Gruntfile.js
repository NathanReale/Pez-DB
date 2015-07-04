module.exports = function(grunt) {

	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		browserify: {
			dist: {
				files: {
					'dist/app.js': ['app/js/app.js']
				},
				options: {
					transform: ['debowerify']
				}
			}
		},
		copy: {
			all: {
				expand: true,
				cwd: 'app/',
				src: ['**/*.html', 'index.html'],
				dest: 'dist/'
			}
		},
		jshint: {
			files: ['Gruntfile.js', 'app/js/*.js', 'app/js/**/*.js'],
			options: {
				globals: {
					jQuery: true,
					console: true,
					document: true
				}
			}
		},
		less: {
			development: {
				options: {
					paths: ['bower_components/bootstrap/less']
				},
				files: {
					"dist/styles.css": "app/styles/style.less"
				}
			}
		},
		watch: {
			js: {
				files: ['<%= jshint.files %>'],
				tasks: ['jshint', 'browserify']
			},
			less: {
				files: ['app/styles/*.less'],
				tasks: ['less']
			},
			copy: {
				files: ['app/**/*.html'],
				tasks: ['copy']
			}
		}
	});

	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-copy');
	grunt.loadNpmTasks('grunt-browserify');

	grunt.registerTask('test', ['jshint']);

	grunt.registerTask('default', ['jshint', 'browserify', 'less', 'copy']);

};
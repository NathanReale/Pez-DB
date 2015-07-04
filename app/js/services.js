angular.module('pezServices', ['ngResource'])
	.factory('Pez', ['$resource', function ($resource) {
		return $resource('api/pez/:id', {id: '@Id'});
	}])
	.factory('Series', ['$resource', function ($resource) {
		return $resource('api/series/:name', {name: '@name'}, {
			get: {method: 'GET', isArray: true}
		});
	}])
	.factory('Wishlist', ['$resource', function ($resource) {
		return $resource('api/wishlist/:id', {id: '@Id'});
	}]);

angular.module('autofillServices', [])
	.service('autofill', ['$http', function($http) {

		var that = this;

		this.categories = [];
		$http.get('/api/categories').then(function (response) {
			that.categories = response.data;
		});

		this.series = [];
		$http.get('/api/series').then(function (response) {
			that.series = response.data;
		});

		this.colors = [];
		$http.get('/api/colors').then(function (response) {
			that.colors = response.data;
		});

		this.feet = ['Yes', 'No', 'Thin'];
		this.originCountry = ['Not Listed', 'China', 'Hungary', 'Slovenia', 'Austria', 'Yugoslavia', 'Czechoslovakia', 'Czech Republic', 'Hong Kong'];
		this.releaseCountry = ['USA', 'Europe', 'Japan', 'Australia'];
		this.patentNumber = ['2.6', '3.4', '3.9', '4.9', '5.9', '7.5', 'None'];
		this.imc = ['1', '2', '3', '4', '5', '6', '7', '8', 'None'];
	}]);
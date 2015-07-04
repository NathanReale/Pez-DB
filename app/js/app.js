require('jquery');

require('angular');
require('angular-route/angular-route');
require('angular-resource/angular-resource');

require('bootstrap/js/alert');
require('bootstrap/js/collapse');
require('angular-bootstrap');
require('angular-strap/dist/angular-strap');

require('./services.js');

var app = angular.module('pez', ['ngRoute', 'pezServices', 'autofillServices', '$strap.directives', 'ui.bootstrap'])
	.config(['$routeProvider', function ($routeProvider) {
		$routeProvider.when('/', {templateUrl: 'dist/partials/home.html'});
		$routeProvider.when('/list', {templateUrl: 'dist/partials/list.html', controller: 'PezCtrl'});
		$routeProvider.when('/series/:name', {templateUrl: 'dist/partials/series.html', controller: 'SeriesCtrl'});
		$routeProvider.when('/add', {templateUrl: 'dist/partials/add.html', controller: 'NewPezCtrl'});
		$routeProvider.when('/edit/:id', {templateUrl: 'dist/partials/add.html', controller: 'EditPezCtrl'});
		$routeProvider.when('/view/:id', {templateUrl: 'dist/partials/view.html', controller: 'PezDetailCtrl'});
		$routeProvider.when('/wishlist', {templateUrl: 'dist/partials/wishlist.html', controller: 'WishlistCtrl'});
		$routeProvider.otherwise({redirectTo: '/'});
	}])
	.run(['$rootScope', function ($rootScope) {
		$rootScope.alerts = [];
		$rootScope.query = '';
	}]);

app.controller('PezCtrl', ['$scope', 'Pez', '$rootScope', function ($scope, Pez, $rootScope) {
	$scope.currentPage = 1;
	$scope.filtered = [];
	$scope.maxSize = 10;
	$scope.itemsPerPage = 20;
	$scope.pez = Pez.query();

	$scope.$watch('filtered.length', function() {
		$scope.totalItems = $scope.filtered.length;
	});

}]);

app.controller('SeriesCtrl', ['$scope', 'Series', '$routeParams', function ($scope, Series, $routeParams) {
	$scope.title = $routeParams.name;
	$scope.pez = Series.get({name: $routeParams.name});
}]);

app.controller('PezDetailCtrl', ['$scope', 'Pez', '$routeParams', '$location', function ($scope, Pez, $routeParams, $location) {
	$scope.data = Pez.get({id: $routeParams.id});

	$scope.edit = function() {
		$location.path("/edit/" + $scope.data.Id);
	};
}]);

app.controller('NewPezCtrl', ['$scope', 'Pez', '$location', '$rootScope', 'autofill',
	function ($scope, Pez, $location, $rootScope, $autofill) {
		$scope.title = "Add Pez";
		$scope.autofill = $autofill;


		$scope.save = function() {
			Pez.save($scope.data, function (data) {
				$rootScope.alerts.push({"type": "success", "content": "New Pez added!"});
				$location.path("/view/" + data.Id);
			});
		};

		$scope.saveMore = function() {
			Pez.save($scope.data, function (data) {
				$rootScope.alerts.push({"type": "success", "content": "New Pez added!"});
				$scope.data.Name = $scope.data.StemColor = $scope.data.Image = $scope.data.Variation =
					$scope.data.Feet = $scope.data.Duplicates = $scope.data.Notes = "";
			});
		};

		$scope.cancel = function() {
			$location.path("/home");
		};
}]);

app.controller('EditPezCtrl', ['$scope', 'Pez', '$routeParams', '$location', '$rootScope', 'autofill',
	function ($scope, Pez, $routeParams, $location, $rootScope, $autofill) {
		$scope.title = "Edit Pez";
		$scope.autofill = $autofill;
		$scope.edit = true;
		$scope.data = Pez.get({id: $routeParams.id});

		$scope.save = function() {
			$scope.data.$save(function (data) {
				$rootScope.alerts.push({"type": "success", "content": "Pez Updated!"});
				$location.path("/view/" + data.Id);
			});
		};

		$scope.remove = function() {
			if (window.confirm("Are you sure you want to delete this Pez?")) {
				$scope.data.$delete(function (data) {
					$location.path("/list");
				});
			}
		};

		$scope.cancel = function() {
			$location.path("/view/" + $scope.data.Id);
		};
	}
]);

app.controller('WishlistCtrl', ['$scope', 'Wishlist', function ($scope, Wishlist) {
	$scope.wishlist = Wishlist.query();

	$scope.showAdd = function() {
		$scope.adding = true;
	};

	$scope.cancelAdd = function() {
		$scope.add = {};
		$scope.adding = false;
	};

	$scope.addNew = function() {
		Wishlist.save($scope.add, function (data) {
			$scope.wishlist = Wishlist.query();
			$scope.add = {};
			$scope.adding = false;
		});
	};

	$scope.remove = function (item) {
		if (window.confirm("Are you sure you want to delete this item from the wishlist?")) {
			item.$delete(function (data) {
				$scope.wishlist = Wishlist.query();
			});
		}
	};

	$scope.edit = function (item) {
		item.editing = true;
	};

	$scope.cancelEdit = function (item) {
		item.editing = false;
	};

	$scope.save = function (item) {
		item.$save(function (data) {
			$scope.wishlist = Wishlist.query();
		});
	};

}]);

app.filter('page', function() {
	return function(input, page) {
		var offset = (page - 1) * 20;
		return input.slice(offset, offset + 20);
	};
});
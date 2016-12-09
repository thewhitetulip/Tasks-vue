/* This is the Vue front end for the Tasks application. We will not be using webpack or any other fancy tech. Plain JS.
 * You are expected to know a little bit of JS. If you are a total newbie, we recommend reading, 
 * https://github.com/getify/You-Dont-Know-JS
 * You certainly do not need to be a pro in JS, just need to know enough to follow along, but  you do need to read the book
 * eventually some day to become comfortable with the concepts in JS, so we recommend strongly to read the book before continuing.
 * */

var app = new Vue({
	//~ this is the element in the html page where Vue will be anchored
	el: '#tasks',
	//~ these are the delimiters which Vue will use, are modified since Go uses {{.
	delimiters: ['${', '}'],
	data: {
		navigation: "pending", //this is what we will display in the title tag of the main page, Completed/Deleted/Pending
		user: "suraj",
		newCategoryName: '', // this is the new category name to be used in the update category flow
		notificationVisible: false, // This toggles the visibility of the notification
		notification: "", // actual content of the notification
		task: {
			ID: "",
			title: "",
			content: "",
			category: '',
			priority: '',
			comments: [],
			showComment: false
		}, // variable in which task value is stored
		selectedCategoryName: '', // if the user has selected the category in the navigation drawer, this var has value
		selectedTaskTypeName: 'pending', // by default we show pending tasks
		comment: {
			content: "",
			author: "",
			created: ""
		}, // data structure to store comment
		category: {
			categoryID: '',
			categoryName: '',
			taskCount: ''
		}, // data structure to store category
		categories: [], // stores all the categories
		tasks: [], // stores all the tasks
		pendingTasks: [],
		categoryTasks: []

	},
	mounted: function () {
		this.fetchTasks();
		this.fetchCategories();
	},
	methods: {
		// This will fetch task from the DB
		fetchTasks: function () {
			var events = [];
			this.$http.get('/api/task/')
				.then(response => response.json())
				.then(result => {
					Vue.set(this.$data, 'tasks', result);
				})
				.catch(err => {
					console.log(err);
				});
		},


		// this will add the task from the user input to our array
		addTask: function (item) {
			this.tasks.push(this.task);
			this.task = {
				title: "",
				content: "",
				category: '',
				priority: '',
				comments: [],
				showComment: false
			}
			$('#addNoteModal').modal('hide');
		},
		// this will add a new category to our data store
		addCategory: function () {
			this.category.taskCount = 0;
			this.categories.push(this.category);
			console.log(this.category.categoryName);
			this.category = {
				categoryID: '',
				categoryName: '',
				taskCount: ''
			};
			this.notificationVisible = true;
			this.notification = "Category Added";
		},
		deleteCategory: function (name) {
			console.log("deleting " + name);
			var index = 0;
			for (category in this.categories) {
				if (this.categories[category].categoryName == name) {
					index = this.categories.indexOf(category);
				}
			}

			this.categories.splice(index, 1);
		},
		// this will add a new note to the existing list of comments
		addComment: function (comment, taskIndex) {
			comment.author = this.user;
			comment.created = new Date();
			if (comment.content != '') {
				this.tasks[taskIndex].comments.push(comment);
				this.comment = {
					content: "",
					created: ""
				}
				this.notification = "added comment";
			} else {
				this.notification = "can't add comment";
			}
			this.notificationVisible = true;
		},
		// will hide the visibility of the notification
		hide: function () {
			this.notificationVisible = false;
		},
		// will delete a comment
		deleteComment: function (taskIndex, commentIndex) {
			this.tasks[taskIndex].comments.splice(commentIndex, 1);
			this.notificationVisible = true;
			this.notification = "Comment deleted";
		},
		// will edit a task
		edit: function (index) {
			t = this.tasks[index];
			this.task.title = t.title;
			this.task.content = t.content;
			this.task.priority = t.priority;
			this.task.category = t.category;
			$('#addNoteModal').modal('show');
		},
		// will delete a task
		deleteTask: function (index) {
			this.tasks.splice(index, 1);
			this.notificationVisible = true;
			this.notification = "Deleted";
		},
		// this will mark the task as completed
		complete: function (index) {
			this.completedTasks.push(this.tasks[index]);
			this.tasks.splice(index, 1);
			this.notificationVisible = true;
			this.notification = "marked as complete";

		},
		// toggles the state to check which part is currently active
		// either pending/complete/deleted or categories
		taskByCategory: function (category) {
			this.selectedCategoryName = category.categoryName;
			this.selectedCategoryID =
				this.selectedTaskTypeName = ''
			this.tasks = [];
		},
		// shows completed tasks
		showCompletedTasks: function (type) {
			this.tasks = this.completedTasks;
			this.selectedTaskTypeName = 'completed';
			this.selectedCategoryName = '';
		},
		// shows pending tasks
		showPendingTasks: function (type) {
			this.tasks = this.pendingTasks;
			this.selectedTaskTypeName = 'pending';
			this.selectedCategoryName = ''
		},
		// shows the deleted tasks
		showDeletedTasks: function (type) {
			this.tasks = this.deletedTasks;
			this.selectedTaskTypeName = 'deleted';
			this.selectedCategoryName = ''
		},
		// used to toggle the visibility of the note's comment area + content area
		toggleContent: function (item) {
			item.showComment = !item.showComment;
		},

		updateCategory: function (oldName, newName) {
			// update the category name in the db
			// this logic is temporary and will be removed later
			var id = '';
			for (category in this.categories) {
				if (this.categories[category].categoryName == oldName) {
					this.categories[category].categoryName = newName;
					console.log("Updated");
				}
			}
		}

	}
})
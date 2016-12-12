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
  delimiters: [
    '${',
    '}'
  ],
  data: {
    navigation: 'Pending', //this is what we will display in the title tag of the main page, Completed/Deleted/Pending
    isLoggedIn: false,
    user: '',
    newCategoryName: '', // this is the new category name to be used in the update category flow
    notificationVisible: false, // This toggles the visibility of the notification
    notification: '', // actual content of the notification
    isEditing: false, // this will decide if we are editing or adding a task, updated whenever we click the edit function, set to true
    taskIDEdit:-1, // stores the task id which is currently being updated;
    categoryEdit: false,
    userLogin : {
    	username: '',
	password:''
    },
    userSignup : {
        username: '',
	password: '',
	email: ''
    },
    task: {
      id: '',
      title: '',
      content: '',
      category: '',
      priority: '',
      comments: [
      ],
      ishidden:'',
      completedmsg:'',
      showComment: false
    }, // variable in which task value is stored
    selectedCategoryName: '', // if the user has selected the category in the navigation drawer, this var has value
    selectedTaskTypeName: 'pending', // by default we show pending tasks
    comment: {
      id:'',
      taskID:'',
      content: '',
      author: '',
      created: ''
    }, // data structure to store comment
    category: {
      categoryID: '',
      categoryName: '',
      taskCount: ''
    }, // data structure to store category
    categories: [
    ], // stores all the categories
    tasks: [
    ], // stores all the tasks
  },
  mounted: function () {
  	this.checklogin();
  },
  methods: {
    signup: function(){
    	this.$http.put('/signup/', this.userSignup, { emulateJSON : true }).then(response => response.json()).then(result => {
		this.notificationVisible=true;
		this.notification = "Sign up successful, pls login"
		this.userSignup = {
		        username: '',
		        password: '',
		        email: ''
		}
	}).catch(err => {
		console.log(err);
	});
    },
    checklogin: function(){
        this.$http.get('/login/').then(response => response.json()).then(result => {
		this.isLoggedIn = result.loggedin;
		this.FetchCategories();
		this.FetchTasks();
	}).catch(err => {
		console.log(err);
	});
    },
    login: function () {
    	this.$http.post('/login/', this.userLogin, { emulateJSON : true }).then(response => response.json()).then(result => {
    		this.isLoggedIn = true;
		this.FetchCategories();
		this.FetchTasks();
		this.userLogin = {username:'', password:''}
	}).catch(err => {
		console.log(err);
	});
    },
    logout: function () {
    	this.$http.get('/logout/').then(response => response.json())
		.then(result => {
			this.isLoggedIn = false;
		}).catch(err => {
			console.log(err);
		});
    },
    // This will fetch task from the DB
    FetchTasks: function () {
      this.tasks = [
      ];
      this.$http.get('/task/').then(response => response.json()).then(result => {
        if (result != null) {
          Vue.set(this.$data, 'tasks', result);
        }
      }).catch (err => {
        console.log(err);
	isLoggedIn = false;
	this.notificationVisible = true;
	this.notification = "unable to fetch tasks";
      });
    },
    FetchCategories: function () {
      this.$http.get('/categories/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'categories', result);
      }).catch (err => {
        console.log(err);
      });
    },
    // this will add the task from the user input to our array
    AddTask: function (item) {
      this.$http.put('/task/', this.task, {
        emulateJSON: true
      }).then(response => response).then(result => {
       if (this.task.ishidden == false) {
           this.tasks.push(this.task);
       }
       categoryIndex = 0;
       for (c in this.categories) {
       	  if (this.categories[c].categoryName== this.task.category) {
	        this.categories[c].taskCount +=1;
	      break;
	  }
       }

        this.task = {
          title: '',
          content: '',
          category: '',
          priority: '',
          comments: [
          ],
          showComment: false
        }
      }).catch (err => {
        console.log(err);
      });
      $('#addNoteModal').modal('hide');
    },
    toggleEditCategoryForm : function(){
    	this.categoryEdit= ! this.categoryEdit;
    },
    // this will add the task from the user input to our array
    UpdateTask: function (item) {
      this.$http.post('/task/', this.task, {
        emulateJSON: true
      }).then(response => response).then(result => {
        // this.tasks.push(this.task);
	index = 0;
	for (t in this.tasks) {
		if (t.id == this.taskIDEdit) {
			index = this.tasks.indexOf(t);
		}
	}
	newTask = this.task;

	this.tasks[index].title = newTask.title;
	this.tasks[index].category = newTask.category;
	this.tasks[index].content = newTask.content;
	this.tasks[index].priority= newTask.priority;

	this.notificationVisible = true;
	this.notification = "Updated task";

        this.task = {
          title: '',
          content: '',
          category: '',
          priority: '',
          comments: [
          ],
          showComment: false
        }
      }).catch (err => {
        console.log(err);
      });
      $('#addNoteModal').modal('hide');
    },
    // this will add a new category to our data store
    addCategory: function () {
      console.log(this.category);
      this.$http.put("/category/", this.category,{
      	emulateJSON: true
      }).then(response => response.json()).then(result => {
      		this.category.taskCount = 0;
      		this.categories.push(this.category);
      		this.category = {
        		categoryID: '',
      			categoryName: '',
        		taskCount: ''
      		};
      		this.notificationVisible = true;
	        this.notification = 'Category Added';
      }).catch (err => {
      		console.log(err);
      		this.notificationVisible = true;
	        this.notification = 'Unable to add category';
      });
    },
    deleteCategory: function (name) {
     this.$http.delete('/category/'+name).then(response => response.json())
     	.then(result => {	
		console.log('deleting ' + name);
		var index = 0;
		for (category in this.categories) {
		   if (this.categories[category].categoryName == name) {
			index = this.categories.indexOf(category);
        	   }
      		}
	      this.categories.splice(index, 1);
	      this.FetchTasks();
	      this.navigation='Pending';
	      this.selectedTaskTypeName = 'pending'
	}).catch(err => {
		console.log(err);
	});
    },
    // this will add a new note to the existing list of comments
    addComment: function (comment, taskIndex) {
     this.comment.taskID = this.tasks[taskIndex].id;
     console.log(this.tasks[taskIndex].title, this.tasks[taskIndex].id);
     this.$http.put('/comment/', this.comment, {
     	emulateJSON: true
     }).then(response => response.json()).then(result => {
        this.comment.author = this.user;
        this.comment.created = new Date();
        if (this.comment.content != '') {
	  if (this.tasks[taskIndex].comments == null) {
	  	this.tasks[taskIndex].comments = [];
	  }
          this.tasks[taskIndex].comments.push(comment);
          this.comment = {
            content : '',
            created: '',
	    taskID:'',
	    author:''
        }
          this.notification = 'added comment';
        } else {
          this.notification = 'can\'t add comment';
        }
      this.notificationVisible = true;
     }).catch(err => {
     	console.log(err);
     });
    },
    // will hide the visibility of the notification
    hide: function () {
      this.notificationVisible = false;
    },
    // will delete a comment
    deleteComment: function (taskIndex, commentIndex, taskID, commentID) {
     this.$http.delete('/comment/'+commentID).then(response => response.json())
     	.then(result => {
              this.tasks[taskIndex].comments.splice(commentIndex, 1);
              this.notificationVisible = true;
              this.notification = 'Comment deleted';	
	}).catch(err => {
	      console.log(err);
	});

    },
    // will edit a task
    edit: function (index) {
      this.isEditing = true;
      t = this.tasks[index];
      this.task.title = t.title;
      this.task.id = t.id;
      this.taskIDEdit = t.id;
      this.task.content = t.content;
      this.task.priority = t.priority;
      this.task.category = t.category;
      $('#addNoteModal').modal('show');
    },
    // will trash a task, won't delete from db
    TrashTask: function (index, taskID) {
      this.$http.delete('/task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(index, 1);
        this.notificationVisible = true;
        this.notification = 'Deleted';
      }).catch(err => {
     	console.log(err); 
        this.notificationVisible = true;
        this.notification = 'Deleted';
      });
    },
    // will restore a task from deleted
    RestoreTask: function (index, taskID) {
      this.$http.get('/restore-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(index, 1);
      })
      this.notificationVisible = true;
      this.notification = 'Restored';
    },
    // this will mark the task as completed
    CompleteTask: function (taskIndex, taskID) {
      this.$http.get('/complete-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(taskIndex, 1);
        console.log('completing ' + taskIndex)
        this.notificationVisible = true;
        this.notification = 'marked as complete';
      }).catch(err => {
      	console.log(err);
        this.notificationVisible = true;
        this.notification = 'unable to marked as complete';
      });
    },
    inComplete: function (taskIndex, taskID) {
      this.$http.get('/incomplete-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(taskIndex, 1);
        console.log('incomplete ' + taskIndex)
      })
      this.notificationVisible = true;
      this.notification = 'marked as incomplete';
    },
    // toggles the state to check which part is currently active
    // either pending/complete/deleted or categories
    taskByCategory: function (category) {
      this.selectedCategoryName = category;
      this.navigation = this.selectedCategoryName;
      this.tasks = [];
      this.selectedTaskTypeName = '';
      this.$http.get('/category/' + this.selectedCategoryName).then(response => response.json()).then(result => {
          if (result!= null) {
            Vue.set(this.$data, 'tasks', result);
          }
      }).catch(err => {
          console.log(err);
	  this.notificationVisible= true;
	  this.notification = "Unable to fetch tasks";
      });
    },
    // shows completed tasks
    showCompletedTasks: function (type) {
      this.tasks = [
      ];
      this.$http.get('/completed/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'tasks', result);
        this.selectedTaskTypeName = 'completed';
        this.navigation = 'Completed';
        this.selectedCategoryName = '';
      }).catch(err => {
        console.log(err);
	this.notificationVisible = true;
	this.notification = "Unable to fetch tasks";
      });
    },
    // shows pending tasks
    showPendingTasks: function (type) {
      this.FetchTasks();
      this.selectedTaskTypeName = 'pending';
      this.navigation = 'Pending';
      this.selectedCategoryName = ''
    },
    // shows the deleted tasks
    showDeletedTasks: function (type) {
      this.tasks = [
      ];
      this.$http.get('/deleted/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'tasks', result)
        this.selectedTaskTypeName = 'deleted';
        this.navigation = 'Deleted';
        this.selectedCategoryName = ''
      }).catch(err => {
     	console.log(err);
	this.notificationVisible = true;
	this.notification = "Unable to fetch tasks";
      });
    },
    // used to toggle the visibility of the note's comment area + content area
    toggleContent: function (item) {
      item.showComment = !item.showComment;
    },
    updateCategory: function (oldName, newName) {
      // update the category name in the db
      // this logic is temporary and will be removed later
      category = {newCategoryName: this.newCategoryName}
      this.$http.post('/category/' + oldName, category, {
      	emulateJSON:true
      }).then(response => response.json()).then(result => {

        for (category in this.categories) {
          if (this.categories[category].categoryName == oldName) {
            this.categories[category].categoryName = newName;
            console.log('Updated');
	    this.navigation = newName;
	    this.toggleEditCategoryForm();
          }
        }
      }).catch(err => {
          console.log(err);
	  this.notificationVisible = true;
	  this.notification = "Unable to update";
      });

    }
  }
})

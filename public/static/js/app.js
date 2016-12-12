/* 
   Author: Suraj Patil http://github.com/thewhitetulip
   License: MIT

   This is the Vue front end for the Tasks application. We will not be using webpack or any other fancy tech. Plain JS.
   You are expected to know a little bit of JS. If you are a total newbie, we recommend reading, 
   https://github.com/getify/You-Dont-Know-JS
   You certainly do not need to be a pro in JS, just need to know enough to follow along, but  you do need to read the book
   eventually some day to become comfortable with the concepts in JS, so we recommend strongly to read the book before continuing.
*/
var app = new Vue({
  // The element in the html page where Vue will be anchored
  el: '#tasks',
  // The delimiters used in our app, standard delimiters modified as Go uses {{.
  delimiters: [
    '${',
    '}'
  ],
  data: {
    navigation: 'Pending', // Displays in the status bar (Completed/Deleted/Pending)
    isLoggedIn: false,
    user: '',
    newCategoryName: '', // The new category name to be used in the update category flow
    notificationVisible: false, // Toggles the visibility of the notification
    notification: '', // Actual content of the notification.
    isEditing: false, // Decides if we are editing or adding a task, updated whenever we click the edit function, set to true
    taskIDEdit: -1, // Stores the task id which is currently being updated;
    categoryEdit: false, // Toggles the display of the update category form
    userLogin: { // For the login form
      username: '',
      password: ''
    },
    userSignup: { // For the signup form
      username: '',
      password: '',
      email: ''
    },
    task: { // For add task form
      id: '',
      title: '',
      content: '',
      category: '',
      priority: '',
      comments: [],
      ishidden: '',
      completedmsg: '',
      showComment: false
    },
    selectedCategoryName: '', // If the user has selected the category in the navigation drawer, this var has value.
    selectedTaskTypeName: 'pending', // If the user is at pending/deleted/completed this var has value.
    comment: { // For the add comment form.
      id: '',
      taskID: '',
      content: '',
      author: '',
      created: ''
    },
    category: { // For add category form.
      categoryID: '',
      categoryName: '',
      taskCount: ''
    },
    categories: [], // Stores all the categories.
    tasks: [], // Stores all the tasks.
  },
  // mounted is called the moment the Vue app is added to the html DOM.
  // when the app is mounted, we check if the user is logged in or not
  // and update the isLoggedIn variable to true if yes, false if not.
  mounted: function () {
    this.checklogin();
  },
  methods: {
    // notify toggles the notification and set the content of the 
    // notification with the argument passed in it.
    notify: function (message) {
      this.notificationVisible = true;
      this.notification = message;
    },
    signup: function () {
      this.$http.put('/signup/', this.userSignup, {
        emulateJSON: true
      }).then(response => response.json()).then(result => {
        this.notify("Sign up successful, pls login");
        this.userSignup = {
          username: '',
          password: '',
          email: ''
        }
      }).catch(err => {
        console.log(err);
        this.notify("Unable to signup");
      });
    },
    // checklogin checks if the user is logged in, if the user is logged in, 
    // fetches tasks and categories.
    checklogin: function () {
      this.$http.get('/login/').then(response => response.json()).then(result => {
        this.isLoggedIn = result.loggedin;
        this.FetchCategories();
        this.FetchTasks();
      }).catch(err => {
        console.log(err);
        this.notify("Something went wrong")
      });
    },
    login: function () {
      this.$http.post('/login/', this.userLogin, {
        emulateJSON: true
      }).then(response => response.json()).then(result => {
        this.isLoggedIn = true;
        this.FetchCategories();
        this.FetchTasks();
        this.userLogin = {
          username: '',
          password: ''
        }
      }).catch(err => {
        console.log(err);
        this.notify("Unable to login");
      });
    },
    logout: function () {
      this.$http.get('/logout/').then(response => response.json())
        .then(result => {
          this.isLoggedIn = false;
        }).catch(err => {
          console.log(err);
          this.notify("Unable to logout");
        });
    },
    FetchTasks: function () {
      this.$http.get('/task/').then(response => response.json()).then(result => {
        if (result != null) {
          Vue.set(this.$data, 'tasks', result);
        } else {
          this.tasks = [];
        }
      }).catch(err => {
        console.log(err);
        isLoggedIn = false;
        this.notify("Unable to fetch Tasks")
      });
    },
    FetchCategories: function () {
      this.$http.get('/categories/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'categories', result);
      }).catch(err => {
        console.log(err);
        this.notify("Unable to fetch categories");
      });
    },
    AddTask: function (item) {
      this.$http.put('/task/', this.task, {
        emulateJSON: true
      }).then(response => response).then(result => {
        if (this.task.ishidden == false) {
          this.tasks.push(this.task);
        }

        this.UpdateCategoryCount(this.task.category, "+", 1);

        this.task = {
          title: '',
          content: '',
          category: '',
          priority: '',
          comments: [],
          showComment: false
        }
      }).catch(err => {
        console.log(err);
        this.notify("Unable to add Task");
      });
      $('#addNoteModal').modal('hide');
    },
    toggleEditCategoryForm: function () {
      this.categoryEdit = !this.categoryEdit;
    },
    // UpdateCategoryCount updates the count of the category where 
    // some action took place, like add task or delete task.
    // our taskCount stores the # of tasks which are pending. 
    // arguments: (category name, action(increase/decrease),value to increase or decrease),
    UpdateCategoryCount: function (name, action, value) {
      categoryIndex = 0;
      for (c in this.categories) {
        if (this.categories[c].categoryName == name) {
          switch (action) {
            case "+": 
                this.categories[c].taskCount += 1;
            case "-":
                this.categories[c].taskCount -= 1;
          }
          break;
        }
      }
    },
    UpdateTask: function (item) {
      this.$http.post('/task/', this.task, {
        emulateJSON: true
      }).then(response => response).then(result => {
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
        this.tasks[index].priority = newTask.priority;

        this.notify("Updated task");

        this.task = {
          title: '',
          content: '',
          category: '',
          priority: '',
          comments: [],
          showComment: false
        }
      }).catch(err => {
        console.log(err);
        this.notify("Unable to update Task");
      });
      $('#addNoteModal').modal('hide');
    },
    addCategory: function () {
      console.log(this.category);
      this.$http.put("/category/", this.category, {
        emulateJSON: true
      }).then(response => response.json()).then(result => {
        this.category.taskCount = 0;
        this.categories.push(this.category);
        this.category = {
          categoryID: '',
          categoryName: '',
          taskCount: ''
        };
        this.notify('Category Added');
      }).catch(err => {
        console.log(err);
        this.notify("Unable to add category");
      });
    },
    deleteCategory: function (name) {
      this.$http.delete('/category/' + name).then(response => response.json())
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
          this.navigation = 'Pending';
          this.selectedTaskTypeName = 'pending'
        }).catch(err => {
          console.log(err);
          this.notify("Unable to delete category");
        });
    },
    addComment: function (comment, taskIndex) {
      this.comment.taskID = this.tasks[taskIndex].id;
      console.log(this.tasks[taskIndex].title, this.tasks[taskIndex].id);

      if (this.comment.content == '') {
        this.notify("Comment can't be empty");
        return;
      }

      this.$http.put('/comment/', this.comment, {
        emulateJSON: true
      }).then(response => response.json()).then(result => {
        this.comment.author = this.user;
        this.comment.created = new Date();

        if (this.tasks[taskIndex].comments == null) {
          this.tasks[taskIndex].comments = [];
        }

        this.tasks[taskIndex].comments.push(comment);
        this.comment = {
          content: '',
          created: '',
          taskID: '',
          author: ''
        }

        this.notify("Comment added")

      }).catch(err => {
        console.log(err);
        this.notify("Unable to add comment");
      });
    },
    // hides the visibility of the notification
    hide: function () {
      this.notificationVisible = false;
    },
    deleteComment: function (taskIndex, commentIndex, taskID, commentID) {
      this.$http.delete('/comment/' + commentID).then(response => response.json())
        .then(result => {
          this.tasks[taskIndex].comments.splice(commentIndex, 1);
          this.notify("Comment deleted");
        }).catch(err => {
          console.log(err);
          this.notify("Unable to delete comment");
        });

    },
    // edit loads a task to the edit form. User can click on the button to update
    // the task, it calls UpdateTask function.
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
    // Trashs a task, won't delete from db.
    TrashTask: function (taskIndex, taskID, category) {
      this.$http.delete('/task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(taskIndex, 1);
        this.notify("Task deleted");
        this.UpdateCategoryCount(category, "-", 1);
      }).catch(err => {
        console.log(err);
        this.notify("Unable to trash Tash");
      });
    },
    RestoreTask: function (taskIndex, taskID, category) {
      this.$http.get('/restore-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(index, 1);
        this.notify("Task restored");
	this.UpdateCategoryCount(category, "+", 1);
      }).catch(err => {
        console.log(err);
        this.notify("Unable to restore task");

      });
    },
    CompleteTask: function (taskIndex, taskID, category) {
      this.$http.get('/complete-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(taskIndex, 1);
        this.notify("Marked task as complete");
	this.UpdateCategoryCount(category, "-", 1);
      }).catch(err => {
        console.log(err);
        this.notify("Unable to mark as complete");
      });
    },
    inComplete: function (taskIndex, taskID, category) {
      this.$http.get('/incomplete-task/' + taskID).then(response => response.json()).then(result => {
        this.tasks.splice(taskIndex, 1);
        this.notify("Marked task as incomplete");
	this.UpdateCategoryCount(category, "+", 1);
      }).catch(err => {
        this.notify("Unable to mark task as incomplete");
      });
    },
    // toggles the state to check which part is currently active
    // either pending/complete/deleted or categories
    taskByCategory: function (category) {
      this.selectedCategoryName = category;
      this.navigation = this.selectedCategoryName;
      this.tasks = [];
      this.selectedTaskTypeName = '';
      this.$http.get('/category/' + this.selectedCategoryName).then(response => response.json()).then(result => {
        if (result != null) {
          Vue.set(this.$data, 'tasks', result);
        }
      }).catch(err => {
        console.log(err);
        this.notify("Unable to fetch tasks");
      });
    },
    showCompletedTasks: function (type) {
      this.$http.get('/completed/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'tasks', result);
        this.selectedTaskTypeName = 'completed';
        this.navigation = 'Completed';
        this.selectedCategoryName = '';
      }).catch(err => {
        console.log(err);
        this.notify("Unable to fetch tasks");
      });
    },
    showPendingTasks: function (type) {
      this.FetchTasks();
      this.selectedTaskTypeName = 'pending';
      this.navigation = 'Pending';
      this.selectedCategoryName = ''
    },
    showDeletedTasks: function (type) {
      this.$http.get('/deleted/').then(response => response.json()).then(result => {
        Vue.set(this.$data, 'tasks', result)
        this.selectedTaskTypeName = 'deleted';
        this.navigation = 'Deleted';
        this.selectedCategoryName = ''
      }).catch(err => {
        console.log(err);
        this.notify("Unable to fetch tasks");
      });
    },
    // Toggles the visibility of the note's comment area + content area
    toggleContent: function (item) {
      item.showComment = !item.showComment;
    },
    updateCategory: function (oldName, newName) {
      category = {
        newCategoryName: this.newCategoryName
      }
      this.$http.post('/category/' + oldName, category, {
        emulateJSON: true
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
        this.notify("Unable to update Task");
      });

    }
  }
})

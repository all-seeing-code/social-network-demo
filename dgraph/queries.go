// Simple queries to add Users, Posts and edges etc.

// Add a user:
mutation {
	addUser(input: {email: "anurag-1@dgraph.io", firstName: "Anurag"}) {
	  user {
		id
		firstName
		email
	  }
	}
}

// Add a post
mutation {
	addPost(input: {title: "My first post", hasCreator: {id: "0x33"}, content: "Dgraph is Awesome"}) {
	  post {
		content
		hasCreator {
		  email
		  firstName
		}
		id
	  }
	}
}


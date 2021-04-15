//Lambda script
// You can type/paste your script here
async function newComment({args, graphql}) {
    // lets add comment first to the db
    const results = await graphql(`mutation ($content: String!, $post: ID!, $user: ID!, $time: Int!) {
        addComment(input: {content: $content, post: {id: $post}, hasCreator: {id: $user}, createdAt: $time }) {
            comment {
                id
                hasCreator {
                    email
                }
            }
        }
    }`, {
        "content": args.content,
        "post": args.postId,
        "user": args.userId,
        "time": args.timestamp,
    })
    console.log("comment created. now sending out email to ", results.data.addComment.comment[0].hasCreator.email)
    return results.data.addComment.comment[0].id
}
self.addGraphQLResolvers({
    "Mutation.newComment": newComment
})
# Blog Service
## About it
- A Blog service needs to the CRUD operations on the blog.

## Requirements
1. As a user, I want to get specific blog information.
2. As a user, I want to create a blog.
3. As a user, I want to delete specific blog information.
4. As a user, I want to update specific blog information.
5. As a user, I want to get a list of blog information.
### Function
1. `func CreateBlog(context.Context, blog) (blog, error)`
    - The function can create a blog based on `blog` information.
2. `func DeleteBlog(context.Context, id, author) (error)`
   - The function can delete a blog based on `id`.
   - The function will determine whether the `author` is the author of the article.
3. `func EditBlog(context.Context, id) (blog, error)`
   - The function can update a blog based on a specific `id`.
4. `func GetBlog(context.Context, id) (blog, error)`
   - The function can get a blog based on a specific `id`.
5. `func GetBlogList(context.Context, pageSize, pageToken, filter) (blogs, nextToken, error)`
   - The function can get a list of blogs based on the `filter`.
   - `pageSize` is the maximum value that how many blogs will give back.
   - `pageToken` is an identifier that is used to get the next page result.
   - `nextToken` is an identifier that represents the next page.
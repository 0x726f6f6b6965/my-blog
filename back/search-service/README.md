# Search Service
## About it
- A search service needs to take some text input, a search query, from the user and return the relevant content in a few seconds or less.

## Requirements
1. As a user, I want to get relevant content based on my search queries.
2. As a user, I want to be able to find the latest content.
### Function
1. `func AutoComplete(ctx context.Context, words string) ([]string, error)`
   - The function runs when a user queries the system to find some content.
   - `words` is the textual query entered by the user in the search bar, based on which the results are found.
2. `func AddIndex(ctx context.Context, index string, id string) (error)`
   - The function runs when content is created.
   - `index` is the topic of content.
   - `id` is the primary key of the database.
3. `func Search(ctx context.Context, query string) ([]string, error)`
   - The function will get the primary keys of the database based on the `query`. 
## Architecture
- Instead of using the database directly, the service will create a Trie structure and inverted indexes in the RAM to implement all the functions we need. With this structure, the response time to a search query will be faster than using the database. 
- The service contains indexing and searching on the same node. Although it seems like efficient usage of resources, it has its downsides as well. Searching and indexing are both resource-intensive operations. Both operations impact the performance of each other. Also, this colocated design doesn’t scale efficiently with varying indexing and search operations over time.

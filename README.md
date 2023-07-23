# beehiveAI
Write a gRPC server for extracting user review data from the PostgreSQL database.
The review format is as follows:
1. Reviewer name. Type - string.
2. Title. Type - string.
3. Text. Type - string.
4. Rating. Type - integer in the range [0, 5].
5. Timestamp. Type - Unix timestamp.
The server should have only one method - `Search()` - and allow the client to filter the data by
rating range (from-to), timestamp (from-to), and reviewer name (full match).
The server should provide high performance for datasets of 100,000+ entries. For testing, you
can use the Amazon reviews(https://cseweb.ucsd.edu/~jmcauley/datasets/amazon/links.html) dataset.
Implementation language - Go.

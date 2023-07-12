
# Election Service API

I created a **REST API** for effective election management using **Go-Lang**, with **PostgresQL** database. Here are some features I have implemented in this API:-

- **Election Officer's Registration** - Election officers can register as the admins for the election service.
- **Election Officer's Sign-in** - Election officers can sign-in into their accounts. The **JWT Authentication** method then provides them with a token, which they can use to access protected routes as shown in the demo.
- **Constituency Registration** - Only authorized election officers can register a new Constituency.
- **Constituency Updation** - Only authorized election officers can manipulate the details of an existing Constituency.
- **Candidate Registration** - Candidates can register themselves to multiple constituencies.
- **Voter Register/Login** - Voters can register and log-in into their accounts.

Here is the demo of working of this API. All the routes are tested using Postman :-


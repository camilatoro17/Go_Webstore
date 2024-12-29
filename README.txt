1. Storing passwords in plain text is dangerous because if anyone accesses the database, they will be able
   to see all the passwords. One solution to this would be to hash the passwords which is basically crypting
   them so that it becomes a string of random characters. There is a good hashing algorithm called bcrypt which
   encrypts the passwords by adding a random string to the password before its hashed so that it's harder to guess.

2. A good way to deal with forgotten passwords is to have the user identify themselves before they can reset it.
   For example, google lets you reset if you authorize it through gmail on your phone, another good way is asking
   for a text code or an email code which verifies that it is the actual owner of the account. A bad way to authorize
   would be to ask security questions because the answers might be easily found on the internet if you really deep
   search into someone.

3. If you implement a remember me, it's important to always give the user the option to sign out and a way to rever
   this decision in case they clicked it by mistake. Another security measure is making sure that this only
   is in effect for a certain amount of time, so this isn't forever maybe just 30 days or less. This is important
   in case the user stops using that device or it's shared with someone else.

4. Create cookies to be secure and use HTTPS only to make sure that they can't be accessed or modified through
   any JavaScript. Making them secure helps them from being manipulated by unencrypted connections. Another best
   practice is to have an expiry date for cookies, not all data needs to be kept forever especially as user
   needs/opinions change and also for security it's best to wipe them after some time.

5. HTTPS is basically HTTP but encrypted. Its main job is making sure that when data is being transferred
   then it stays completely safe and can't be intercepted. It especially protects passwords and credit card
   information which is usually what hackers are looking for. It protects against these by using encryption
   methods in order to protect the data.
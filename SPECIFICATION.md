# SOFTWARE SPECIFICATION

## PURPOSE
This software intended to manage user authentication and authorization.

## SCOPE
This software only related to user's entity, specifically in managing user's session.

* User can register with their name, username and password. Username must be unique, so user can login only by remembering and knowing their username and password.
* User can login with their registered username and password.
* After register or login, user can get their token to make request in protected end-point.
* User can see their profile using their token.


## Limitation
For it's simplicity, in this first phase, this software is not including:
* Change the user password
* Reset the user password when user forget their password
* Sending the email to verify user. This software is not save any email address.


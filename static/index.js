function authStateCheck() {
    firebase.auth().onAuthStateChanged(function (user) {
        if (user) {
            // User is signed in.
            document.getElementById('sign-out').hidden = false;
            document.getElementById('post-form').hidden = false;
            document.getElementById('account-details').textContent =
                'Signed in as ' + user.displayName + ' (' + user.email + ')';
            user.getIdToken().then(function (accessToken) {
                // Add the token to the post form. The user info will be extracted
                // from the token by the server.
                document.getElementById('token').value = accessToken;
                alert('Signed in as ' + user.displayName + ' (' + user.email + ')')
            });
        } else {
            // User is signed out.
            // Initialize the FirebaseUI Widget using Firebase.
            window.location.replace("/");
            alert('User is not sign in')
        }
    }, function (error) {
        console.log(error);
        alert('Unable to log in: ' + error)
    });
    // [END auth_request]
}

function handleSignUp() {
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;
    if (email.length < 4) {
        alert('Please enter an email address.');
        return;
    }
    if (password.length < 4) {
        alert('Please enter a password.');
        return;
    }
    // Sign in with email and pass.
    // [START createwithemail]
    firebase.auth().createUserWithEmailAndPassword(email, password).catch(function(error) {
        // Handle Errors here.
        var errorCode = error.code;
        var errorMessage = error.message;
        // [START_EXCLUDE]
        if (errorCode == 'auth/weak-password') {
            alert('The password is too weak.');
        } else {
            alert(errorMessage);
        }
        console.log(error);
        // [END_EXCLUDE]
    });
    // [END createwithemail]
    document.getElementById("post-form").submit();
    alert('created account and posted form')
}

/**
 * Handles the sign in button press.
 */
function handleSignIn() {
    if (firebase.auth().currentUser) {
        // [START signout]
        firebase.auth().signOut();
        // [END signout]
    } else {
        var email = document.getElementById('email').value;
        var password = document.getElementById('password').value;
        if (email.length < 4) {
            alert('Please enter an email address.');
            return;
        }
        if (password.length < 4) {
            alert('Please enter a password.');
            return;
        }
        // Sign in with email and pass.
        // [START authwithemail]
        firebase.auth().signInWithEmailAndPassword(email, password).catch(function(error) {
            // Handle Errors here.
            var errorCode = error.code;
            var errorMessage = error.message;
            // [START_EXCLUDE]
            if (errorCode === 'auth/wrong-password') {
                alert('Wrong password.');
            } else {
                alert(errorMessage);
            }
            console.log(error);
            // [END_EXCLUDE]
        });
        // [END authwithemail]
    }
    document.getElementById("post-form").submit();
    alert('Handle sign in and post form done')
}
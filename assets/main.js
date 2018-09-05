document.addEventListener('DOMContentLoaded', () => {

    // Get all "navbar-burger" elements
    const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

    // Check if there are any navbar burgers
    if ($navbarBurgers.length > 0) {

        // Add a click event on each of them
        $navbarBurgers.forEach(el => {
            el.addEventListener('click', () => {

                // Get the target from the "data-target" attribute
                const target = el.dataset.target;
                const $target = document.getElementById(target);

                // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
                el.classList.toggle('is-active');
                $target.classList.toggle('is-active');

            });
        });
    }

    const deleteButtons = document.getElementsByClassName('delete');

    for (var i = 0; i < deleteButtons.length; i++) {
        deleteButtons[i].addEventListener('click', dismiss);
    }

    function dismiss(e) {
        this.parentNode.classList.add('is-hidden');
    }

});

function charcountupdate(str) {
    var lng = str.length;
    document.getElementById("charcount").innerHTML = lng + '/200';
}

function toggleShowPassword() {
    var x = document.getElementById("password");
    var y = document.getElementById("showPassword")
    if (x.type === "password") {
        x.type = "text";
    } else {
        x.type = "password";
    }
    if (y.textContent === "Show") {
        y.textContent = "Hide"
    } else {
        y.textContent = "Show"
    }
}

// // delete-warning-modal
// var deleteWarningModal = document.getElementById("delete-warning-modal");
// // open-warning-modal
// var openWarningModalButton = document.getElementById("open-warning-modal");
// //close-warning-modal
// var closeWarningModalButton = document.getElementById("close-warning-modal");

// openWarningModalButton.onclick = function () {
//     deleteWarningModal.classList.add('is-active')
//     console.log(deleteWarningModal.classList)
// }
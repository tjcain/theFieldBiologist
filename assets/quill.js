var quill = new Quill('#editor-container', {
    modules: {
        toolbar: [
            [{
                'header': [1, 2, 3, 4, 5, 6, false]
            }],
            ['bold', 'italic', 'underline', 'strike'], // toggled buttons
            ['blockquote', 'code-block'],
            ['image'],
            [{
                'header': 1
            }, {
                'header': 2
            }], // custom button values
            [{
                    'list': 'ordered'
                },
                {
                    'list': 'bullet'
                }
            ],
            [{
                'script': 'sub'
            }, {
                'script': 'super'
            }], // superscript/subscript
            [{
                'color': []
            }, {
                'background': []
            }], // dropdown with defaults from theme
            [{
                'align': []
            }],
            ['clean'] // remove formatting button
        ]
    },
    placeholder: 'Inform the wrorld...',
    theme: 'snow'
});
var form = document.querySelector('#wysiwygform');
form.onsubmit = function () {
    // Populate hidden form on submit
    var article = document.querySelector('input[name=body]');
    // article.value = JSON.stringify(quill.getContents());
    article.value = quill.container.querySelector('.ql-editor').innerHTML
    // this poops out shitty js object
    console.log("Submitted", $(form).serialize(), $(form).serializeArray());
    // this poops out html
    console.log(quill.container.querySelector('.ql-editor').innerHTML)
    form.submit()
    console.log("subimitted")
    return false;
};
// function sendForm() {
//     form.submit()
//     console.log("subimitted")
// }
// A custom image handler could look like:
//     var imageHandler = function (image, callback) {
//         var formData = new FormData();
//         formData.append('image', image, image.name);
//         var xhr = new XMLHttpRequest();
//         xhr.open('POST', 'handler.php', true);
//         xhr.onload = function () {
//             if (xhr.status === 200) {
//                 callback(xhr.responseText);
//             }
//         };
//         xhr.send(formData);
//     };
// Then Add it to quill like so:
// var quill = new Quill('#editor-container', {
//     ...,
//     imageHandler: imageHandler
// });
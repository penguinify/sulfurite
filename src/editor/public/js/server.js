
urlParams = new URLSearchParams(window.location.search);
title = document.getElementById('filename');

if (urlParams.has('file')) {
    var file = urlParams.get('file');

    title.innerText = file;

    fetch(window.location.origin + '/api/file?file=' + file).then(response => {
        if (response.ok) {
            response.text().then(text => {
                var editor = new Editor(document.getElementById('editor'), text.split('\n').filter(line => line.length > 0));
                editor.Initialize();
                editor.Render();
            });
        } else {
            alert('Error: ' + response.status);
            console.error(response);
        }
    })
}

setInterval(() => {
    if (document.hidden) return;

    // if tab is open check if the server is still running
    fetch(window.location.origin + '/api/heartbeat').catch(() => {
        window.location.reload();
    })
}, 5000);

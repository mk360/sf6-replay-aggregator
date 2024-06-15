document.getElementById("ui").onchange = function(e) {
    const { checked } = e.target;
    if (checked) {
        document.body.classList.add("enable-ui");
    } else {
        document.body.classList.remove("enable-ui");
    }
}
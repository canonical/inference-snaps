// if we already have an onload function, save that one
var prev_handler = window.onload;

window.onload = function onLoadHandler() {
    // Prevent self-recursion
    if (prev_handler && prev_handler !== onLoadHandler) {
        prev_handler();
    }

    const link = document.createElement("a");
    link.classList.add("muted-link");
    link.classList.add("github-issue-link");
    link.text = "Give feedback";
    link.href = (
        github_url
        + "/issues/new/choose"
    );
    link.target = "_blank";

    const div = document.createElement("div");
    div.classList.add("github-issue-link-container");
    div.append(link)

    const container = document.querySelector(".article-container > .content-icon-container");
    container.prepend(div);
};

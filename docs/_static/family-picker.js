// Labels each family panel with its family name, shown when the "All" pill is selected.
// Opens the family tab that contains the URL fragment target (for example #gemma3)
// when that snap is not already visible, so links to any snap keep working.
document.addEventListener("DOMContentLoaded", () => {
  for (const panel of document.querySelectorAll(".family-picker > .sd-tab-content:not(.all-snaps)")) {
    panel.dataset.family = panel.previousElementSibling.textContent.trim();
  }
  openFamilyTabForHash();
});

function openFamilyTabForHash() {
  const id = decodeURIComponent(window.location.hash.slice(1));
  const target = id && document.getElementById(id);
  const panel = target && target.closest(".family-picker > .sd-tab-content");
  if (!panel) return;
  if (!target.offsetParent) {
    const label = panel.previousElementSibling;
    const input = label && document.getElementById(label.htmlFor);
    if (input) input.checked = true;
  }
  target.scrollIntoView();
}

window.addEventListener("hashchange", openFamilyTabForHash);

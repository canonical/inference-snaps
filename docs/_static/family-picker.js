// Opens the family tab that contains the URL fragment target (for example #gemma3),
// so links to a snap in a hidden tab keep working.
function openFamilyTabForHash() {
  const id = decodeURIComponent(window.location.hash.slice(1));
  const target = id && document.getElementById(id);
  const panel = target && target.closest(".family-picker > .sd-tab-content");
  if (!panel) return;
  const label = panel.previousElementSibling;
  const input = label && document.getElementById(label.htmlFor);
  if (input) input.checked = true;
  target.scrollIntoView();
}

document.addEventListener("DOMContentLoaded", openFamilyTabForHash);
window.addEventListener("hashchange", openFamilyTabForHash);

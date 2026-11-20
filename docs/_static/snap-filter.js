// Adds a search box above the family picker on the Available snaps page.
// Typing shows matching snaps from all families; clearing it restores the tabs.
// Without JavaScript the page works as before, with no search box.
document.addEventListener("DOMContentLoaded", () => {
  const picker = document.querySelector(".family-picker");
  if (!picker) return;

  const normalize = (text) => text.toLowerCase().replace(/\s+/g, " ").trim();
  const compact = (text) => text.toLowerCase().replace(/[^a-z0-9]/g, "");

  // Index each card once: its name, family and snap name, plus all its text
  const panels = [...picker.querySelectorAll(":scope > .sd-tab-content")];
  const cards = [];
  for (const panel of panels) {
    const family = panel.previousElementSibling.textContent.trim();
    panel.dataset.family = family;
    for (const card of panel.querySelectorAll(".snap-card")) {
      const name = card.querySelector(".snap-name")?.textContent ?? "";
      const snap = card.querySelector('img[alt$=" snap"]')?.alt.replace(/ snap$/, "") ?? "";
      const title = `${name} ${family} ${snap}`;
      cards.push({
        card,
        panel,
        title: normalize(title),
        titleCompact: compact(title),
        text: normalize(`${title} ${card.textContent}`),
      });
    }
  }

  // A card matches if the whole query appears anywhere in it ("amd gpu", "jetson"),
  // or if every word matches its name, family or snap name ("nemotron nano", "qwen35")
  function matches(entry, query) {
    if (entry.text.includes(query)) return true;
    return query.split(" ").every(
      (word) => entry.title.includes(word) || entry.titleCompact.includes(compact(word)),
    );
  }

  const box = document.createElement("div");
  box.className = "snap-filter";
  box.innerHTML = `
    <label class="snap-filter-label" for="snap-filter-input">Search snaps</label>
    <input id="snap-filter-input" type="search" autocomplete="off" spellcheck="false"
      placeholder="Model, family or hardware, for example qwen, ocr, jetson or amd gpu">
    <p class="snap-filter-status" role="status" aria-live="polite"></p>`;
  picker.before(box);
  const input = box.querySelector("input");
  const status = box.querySelector(".snap-filter-status");

  function applyFilter() {
    const query = normalize(input.value);
    const filtering = query.length > 0;
    picker.classList.toggle("is-filtering", filtering);

    let count = 0;
    for (const entry of cards) {
      const show = !filtering || matches(entry, query);
      entry.card.classList.toggle("snap-filter-hidden", !show);
      if (show) count++;
    }
    for (const panel of panels) {
      const hasMatch = cards.some((e) => e.panel === panel && !e.card.classList.contains("snap-filter-hidden"));
      panel.classList.toggle("snap-filter-hidden", filtering && !hasMatch);
    }

    status.textContent = !filtering
      ? ""
      : count === 0
        ? `No snaps match "${input.value.trim()}".`
        : `${count} ${count === 1 ? "snap matches" : "snaps match"} "${input.value.trim()}".`;

    // Keep the query in the URL so a filtered view can be shared
    const url = new URL(window.location);
    if (filtering) url.searchParams.set("q", input.value.trim());
    else url.searchParams.delete("q");
    history.replaceState(null, "", url);
  }

  input.addEventListener("input", applyFilter);
  input.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      input.value = "";
      applyFilter();
    }
  });

  const initial = new URLSearchParams(window.location.search).get("q");
  if (initial) {
    input.value = initial;
    applyFilter();
  }
});
